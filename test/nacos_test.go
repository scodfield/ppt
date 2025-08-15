package test

import (
	"fmt"
	"go.uber.org/zap"
	"ppt/log"
	"ppt/nacos/wrapper"
	"testing"
)

func TestNacosClient(t *testing.T) {
	dbConf, err := wrapper.GetNacosDBConfig()
	if err != nil {
		log.Error("GetNacosDBConfig error", zap.Error(err))
		return
	}
	log.Info("nacos_db_conf: ", zap.Any("nacos_db_conf", dbConf))
	fmt.Printf("dbConf: %+v\n", dbConf)
}
