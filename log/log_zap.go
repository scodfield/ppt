package log

import "go.uber.org/zap"

func Info(format string, fields ...zap.Field) {
	if Logger.log != nil {
		if fields == nil {
			Logger.log.Info(format)
		} else {
			Logger.log.Info(format, fields...)
		}
	}
}

func Warn(format string, fields ...zap.Field) {
	if Logger.log != nil {
		if fields == nil {
			Logger.log.Warn(format)
		} else {
			Logger.log.Warn(format, fields...)
		}
	}
}

func Error(format string, fields ...zap.Field) {
	if Logger.log != nil {
		if fields == nil {
			Logger.log.Error(format)
		} else {
			Logger.log.Error(format, fields...)
		}
	}
}

func Fatal(format string, fields ...zap.Field) {
	if Logger.log != nil {
		if fields == nil {
			Logger.log.Fatal(format)
		} else {
			Logger.log.Fatal(format, fields...)
		}
	}
}

func Panic(format string, fields ...zap.Field) {
	if Logger.log != nil {
		if fields == nil {
			Logger.log.Panic(format)
		} else {
			Logger.log.Panic(format, fields...)
		}
	}
}
