package config

type ESConfig struct {
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	APIKey   string `json:"api_key"`
}

type CKConfig struct {
	Url      string `json:"url"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Database string `json:"database"`
}
