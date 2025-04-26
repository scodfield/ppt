package config

import "os"

var (
	TimeZone string
	AppName  string
	Env      string
	HostName string
)

func InitGlobalConfig() {
	TimeZone = os.Getenv("TZ")
	AppName = os.Getenv("APP_NAME")
	Env = os.Getenv("ENV")
	HostName, _ = os.Hostname()
}
