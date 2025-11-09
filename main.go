package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/judwhite/go-svc"
	"go.uber.org/zap"
	"net/http"
	pptCache "ppt/cache"
	"ppt/config"
	"ppt/dao"
	"ppt/kafka"
	"ppt/log"
	"ppt/monitor"
	"ppt/mq"
	"ppt/nacos/wrapper"
	"ppt/router"
	"ppt/timer"
	"runtime/debug"
	"sync"
	"syscall"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (wg *WaitGroupWrapper) Wrap(cb func()) {
	wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("WaitGroupWrapper recover panic", zap.Any("recover_err", err))
				debug.PrintStack()
			}
		}()
		cb()
	}()
}

type program struct {
	WaitGroupWrapper
	httpServer *router.HttpServer
	port       int
}

func (s *program) Init(env svc.Environment) error {
	s.initPort()
	monitor.InitProm()
	config.InitGlobalConfig()

	var err error
	err = log.InitUberZap()
	if err != nil {
		fmt.Println("InitUberZap error", zap.Error(err))
		return err
	}

	dbCfg, err := wrapper.GetNacosDBConfig()
	if err != nil {
		log.Error("GetNacosDBConfig error", zap.Error(err))
		return err
	}

	if err = dao.InitRedis(&dbCfg.RedisConfig); err != nil {
		log.Error("ppt init redis error", zap.Error(err))
		return err
	}

	if err = dao.InitPg(&dbCfg.PgConfig); err != nil {
		log.Error("ppt init pg error", zap.Error(err))
		return err
	}

	if err = dao.InitMongo(dbCfg); err != nil {
		log.Error("ppt init mongo error", zap.Error(err))
		return err
	}

	if err = pptCache.InitUserCache(); err != nil {
		log.Error("ppt cache init user error", zap.Error(err))
		return err
	}

	if err = kafka.InitKafkaSarama(&dbCfg.KafkaConfig); err != nil {
		log.Error("ppt init kafka error", zap.Error(err))
		return err
	}

	if err = mq.InitAsynqServer(&dbCfg.RedisConfig); err != nil {
		log.Error("ppt init asynq server error", zap.Error(err))
		return err
	}

	if err = mq.InitAsynq(&dbCfg.RedisConfig); err != nil {
		log.Error("ppt init asynq client error", zap.Error(err))
		return err
	}

	if err = timer.InitTimer(); err != nil {
		log.Error("ppt init timer error", zap.Error(err))
		return err
	}

	s.httpServer = router.NewHttpServer(s.port)

	return nil
}

func (s *program) Start() error {
	s.Wrap(func() {
		mq.StartAsynqServer()
	})

	s.Wrap(func() {
		kafka.StartSaramaKafka()
	})

	s.Wrap(func() {
		if err := s.httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("ppt start http server error", zap.Error(err))
			panic(err)
		}
	})

	log.Info("ppt start http server success")
	return nil
}

func (s *program) Stop() error {
	s.httpServer.Stop()
	mq.CloseAsynq()
	mq.CloseAsynqServer()

	dao.CloseRedis()
	dao.ClosePg()
	dao.CloseMongo()
	kafka.CloseSaramaKafka()
	return nil
}

func (s *program) initPort() {
	port := flag.Int("port", 8090, "Http port number")
	flag.Parse()
	s.port = *port
}

func main() {
	app := &program{}
	if err := svc.Run(app, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT); err != nil {
		log.Error("ppt main run error", zap.Error(err))
	}
}
