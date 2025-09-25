package config

import (
	"os"
	"strconv"
)

var (
	TimeZone  string
	AppName   string
	Env       string
	HostName  string
	NacosHost string
	NacosPort int
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func InitGlobalConfig() {
	TimeZone = os.Getenv("TZ")
	AppName = os.Getenv("APP_NAME")
	Env = os.Getenv("ENV")
	HostName, _ = os.Hostname()
	NacosHost = os.Getenv("NACOS_HOST")
	NacosPort, _ = strconv.Atoi(os.Getenv("NACOS_PORT"))
}
