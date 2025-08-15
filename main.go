package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/judwhite/go-svc"
	"go.uber.org/zap"
	"net/http"
	pptCache "ppt/cache"
	"ppt/dao"
	"ppt/log"
	"ppt/monitor"
	"ppt/router"
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

	var err error
	err = log.InitUberZap()
	if err != nil {
		fmt.Println("InitUberZap error", zap.Error(err))
		return err
	}
	redisCfg := &dao.RedisConfig{}
	if err = dao.InitRedis(redisCfg); err != nil {
		log.Error("ppt init redis error", zap.Error(err))
		return err
	}

	pgCfg := &dao.PgConfig{}
	if err = dao.InitPg(pgCfg); err != nil {
		log.Error("ppt init pg error", zap.Error(err))
		return err
	}

	mongoCfg := &dao.MongoConfig{}
	if err = dao.InitMongo(mongoCfg); err != nil {
		log.Error("ppt init mongo error", zap.Error(err))
		return err
	}

	pptCache.InitUserCache(dao.PgDB, dao.RedisDB, dao.UserCacheDefaultExpiration, dao.UserCacheDefaultCleanUp)
	s.httpServer = router.NewHttpServer(s.port)

	return nil
}

func (s *program) Start() error {
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

	dao.CloseRedis()
	dao.ClosePg()
	dao.CloseMongo()
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
