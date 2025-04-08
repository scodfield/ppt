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
	
	return nil
}
