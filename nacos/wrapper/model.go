package wrapper

type DBConfig struct {
	RedisConfig RedisConfig `json:"redis"`
	PgConfig    PgConfig    `json:"pg"`
	MongoConfig string      `json:"mongo"`
}

type RedisConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	UserName    string `json:"user_name,omitempty"`
	Password    string `json:"password"`
	IsCluster   bool   `json:"is_cluster,omitempty"`
	DBIndex     int    `json:"db,omitempty"`
	SSLVerify   bool   `json:"ssl_verify,omitempty"`
	SSLCaCert   string `json:"ssl_ca_cert,omitempty"`
	SSLCertfile string `json:"ssl_cert_file,omitempty"`
	SSLKeyfile  string `json:"ssl_key_file,omitempty"`
}

type PgConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSLMode  string `json:"ssl_mode,omitempty"`
}
