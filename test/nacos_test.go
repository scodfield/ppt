package test

import (
	"fmt"
	"go.uber.org/zap"
	"ppt/logger"
	"ppt/nacos/wrapper"
	"testing"
)

func TestNacosClient(t *testing.T) {
	dbConf, err := wrapper.GetNacosDBConfig()
	if err != nil {
		logger.Error("GetNacosDBConfig error", zap.Error(err))
		return
	}
	fmt.Printf("dbConf: %+v\n", dbConf)
}
