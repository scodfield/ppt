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

// RedisConfig Redis初始化配置
type RedisConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	UserName    string `json:"user_name,omitempty"`
	Password    string `json:"password"`
	IsClustered bool   `json:"is_clustered,omitempty"`
	DBIndex     int    `json:"db,omitempty"`
	SSLVerify   bool   `json:"ssl_verify"`
	SSLCaCerts  string `json:"ssl_ca_certs,omitempty"`
	SSLCertfile string `json:"ssl_certfile,omitempty"`
	SSLKeyfile  string `json:"ssl_keyfile,omitempty"`
}
