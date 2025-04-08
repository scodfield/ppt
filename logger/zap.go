package logger

import "go.uber.org/zap"

var Logger *zap.Logger

func InitUberZap() error {
	logger, err := zap.NewProduction()
	if err != nil || logger == nil {
		return err
	}
	defer logger.Sync()

	Logger = logger

	// 底层 cores.Write/2 方法同步写;
	// 若自定义 core.Write/2 方法为异步写, 则在调用 buf.Free/0 方法前, 显式调用 core.Sync/0 方法 flush buffer;
	// 或 make([]byte, len(buf.bs), len(buf.bs)), 将 buf 内容copy出来, 以免调用 buf.Free/0 释放 buf 后其它
	// 协程持续引用该buf, 导致被后续日志覆写;
	Logger.Info("init zap logger success")

	return nil
}
