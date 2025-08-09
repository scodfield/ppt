package main

import (
	"flag"
	"github.com/judwhite/go-svc"
	"go.uber.org/zap"
	"ppt/dao"
	"ppt/logger"
	"ppt/router"
	"runtime/debug"
	"sync"
	"syscall"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (wg *WaitGroupWrapper) AddDone(cb func()) {
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

	var err error
	redisCfg := &dao.RedisConfig{}
	if err = dao.InitRedis(redisCfg); err != nil {
		logger.Error("ppt init redis error", zap.Error(err))
		return err
	}

	return nil
}

func (s *program) Start() error {
	return nil
}

func (s *program) Stop() error {
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
