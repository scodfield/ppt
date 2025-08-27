package wrapper

type DBConfig struct {
	RedisConfig RedisConfig `json:"redis"`
	PgConfig    PgConfig    `json:"pg"`
	MongoConfig string      `json:"mongo"`
	KafkaConfig KafkaConfig `json:"kafka"`
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

type KafkaConfig struct {
	BootstrapServer string `json:"bootstrap_servers"`
	Retries         int    `json:"retries"`
	RetryBackoffMs  int    `json:"retry_backoff_ms"`
	LingerMs        int    `json:"linger_ms"`
	Partitions      int    `json:"partitions"`
	BatchSize       int    `json:"batch_size"`
	BufferMemory    int64  `json:"buffer_memory"`
	MaxBlockMs      int    `json:"max_block_ms"`
}
