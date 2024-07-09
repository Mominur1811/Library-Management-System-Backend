package config

type DB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	DbName   string `json:"dbName"`
	SSLMode  string `json:"sslMode"`
	Password string `json:"password"`
	User     string `json:"user"`
}

type Mode string
type Config struct {
	Mode         Mode   `json:"mode"`
	ServiceName  string `json:"service_name"`
	HttpPort     int    `json:"http_port"`
	JwtSecretKey string `json:"jwt_secrect_key"`
	Db           DB     `json:"db"`
}

var config *Config

func init() {
	config = &Config{}
}

func GetConfig() Config {
	return *config
}
