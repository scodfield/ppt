package dao

// PgConfig PostgreSql初始化配置
type PgConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
	Database string `json:"database"`
}

// MongoConfig MongoDB初始化配置
type MongoConfig struct {
	Url                string `json:"url"`
	SecondaryPreferred bool   `json:"secondary_preferred"`
}
