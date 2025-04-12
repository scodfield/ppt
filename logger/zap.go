package logger

import (
	"fmt"
	"go.uber.org/zap"
)

var Logger LoggerV2

func InitUberZap() error {
	logger, err := zap.NewProduction()
	if err != nil || logger == nil {
		return err
	}
	defer logger.Sync()

	Logger.log = logger

	// 底层 cores.Write/2 方法同步写;
	// 若自定义 core.Write/2 方法为异步写, 则在调用 buf.Free/0 方法前, 显式调用 core.Sync/0 方法 flush buffer;
	// 或 make([]byte, len(buf.bs), len(buf.bs)), 将 buf 内容copy出来, 以免调用 buf.Free/0 释放 buf 后其它
	// 协程持续引用该buf, 导致被后续日志覆写;
	Logger.log.Info("init zap logger success")

	return nil
}

func FormatLogV2(f interface{}, v ...interface{}) (placeholder string, fields []zap.Field) {
	switch f.(type) {
	case string:
		placeholder = f.(string)
	default:
		placeholder = fmt.Sprint(f)
	}
	if len(v) <= 0 {
		fields = nil
		return
	}

	fields = make([]zap.Field, 0, len(v))
	for i, k := range v {
		paramName := fmt.Sprintf("param_%d", i)
		switch k.(type) {
		case string:
			fields = append(fields, zap.String(paramName, k.(string)))
		case error:
			fields = append(fields, zap.NamedError(paramName, k.(error)))
		case int:
			fields = append(fields, zap.Int(paramName, k.(int)))
		case int8:
			fields = append(fields, zap.Int8(paramName, k.(int8)))
		case int16:
			fields = append(fields, zap.Int16(paramName, k.(int16)))
		case int32:
			fields = append(fields, zap.Int32(paramName, k.(int32)))
		case zap.Field:
			fields = append(fields, k.(zap.Field))
		case *zap.Field:
			fields = append(fields, *k.(*zap.Field))
		default:
			fields = append(fields, zap.Any(paramName, k))
		}
	}
	return
}

type LoggerV2 struct {
	log *zap.Logger
}

func (l *LoggerV2) Debug(f string, args ...interface{}) {
	placeHolder, fields := FormatLogV2(f, args...)
	if fields == nil {
		l.log.Debug(placeHolder)
	} else {
		l.log.Debug(placeHolder, fields...)
	}
}

func (l *LoggerV2) Info(f string, args ...interface{}) {
	placeHolder, fields := FormatLogV2(f, args...)
	if fields == nil {
		l.log.Info(placeHolder)
	} else {
		l.log.Info(placeHolder, fields...)
	}
}

func (l *LoggerV2) Warn(f string, args ...interface{}) {
	placeHolder, fields := FormatLogV2(f, args...)
	if fields == nil {
		l.log.Warn(placeHolder)
	} else {
		l.log.Warn(placeHolder, fields...)
	}
}

func (l *LoggerV2) Error(f string, args ...interface{}) {
	placeHolder, fields := FormatLogV2(f, args...)
	if fields == nil {
		l.log.Error(placeHolder)
	} else {
		l.log.Error(placeHolder, fields...)
	}
}

func (l *LoggerV2) Fatal(f string, args ...interface{}) {
	placeHolder, fields := FormatLogV2(f, args...)
	if fields == nil {
		l.log.Fatal(placeHolder)
	} else {
		l.log.Fatal(placeHolder, fields...)
	}
}

func (l *LoggerV2) Panic(f string, args ...interface{}) {
	placeHolder, fields := FormatLogV2(f, args...)
	if fields == nil {
		l.log.Panic(placeHolder)
	} else {
		l.log.Panic(placeHolder, fields...)
	}
}

func (l *LoggerV2) Debugf(format string, args ...interface{}) {
	l.Debug(format, args...)
}

func (l *LoggerV2) Infof(format string, args ...interface{}) {
	l.Info(format, args...)
}

func (l *LoggerV2) Warnf(format string, args ...interface{}) {
	l.Warn(format, args...)
}

func (l *LoggerV2) Errorf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}

func (l *LoggerV2) Fatalf(format string, args ...interface{}) {
	l.Fatal(format, args...)
}

func (l *LoggerV2) Panicf(format string, args ...interface{}) {
	l.Panic(format, args...)
}

func (l *LoggerV2) SetLevel(lv int) {

}
