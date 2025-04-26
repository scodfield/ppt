package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitAppConfig() {
	wd, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(wd + "/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func GetHttpPort() string {
	return viper.GetString("server.port")
}
