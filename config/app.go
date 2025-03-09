package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var (
	AppName *string
	Env     *string
)

func InitAppConfig() {
	wd, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(wd + "/config")

	AppName = flag.String("appName", "ppt", "app name")
	Env = flag.String("env", "test", "app run env")
	flag.Parse()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func GetHttpPort() string {
	return viper.GetString("server.port")
}
