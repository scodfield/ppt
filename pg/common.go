package pg

import "context"

type PgConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
	Database string `json:"database"`
}

var (
	Ctx = context.Background()
)
