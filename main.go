package main

import (
	"errors"
	"flag"
	"github.com/judwhite/go-svc"
	"go.uber.org/zap"
	"net/http"
	"ppt/dao"
	"ppt/logger"
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
				logger.Error("WaitGroupWrapper recover panic", zap.Any("recover_err", err))
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
	redisCfg := &dao.RedisConfig{}
	if err = dao.InitRedis(redisCfg); err != nil {
		logger.Error("ppt init redis error", zap.Error(err))
		return err
	}

	pgCfg := &dao.PgConfig{}
	if err = dao.InitPg(pgCfg); err != nil {
		logger.Error("ppt init pg error", zap.Error(err))
		return err
	}

	mongoCfg := &dao.MongoConfig{}
	if err = dao.InitMongo(mongoCfg); err != nil {
		logger.Error("ppt init mongo error", zap.Error(err))
		return err
	}

	s.httpServer = router.NewHttpServer(s.port)

	return nil
}

func (s *program) Start() error {
	s.Wrap(func() {
		if err := s.httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("ppt start http server error", zap.Error(err))
			panic(err)
		}
	})

	logger.Info("ppt start http server success")
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
		logger.Error("ppt main run error", zap.Error(err))
	}
}
