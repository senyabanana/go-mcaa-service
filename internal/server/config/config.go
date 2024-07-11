package config

type Config struct {
	Address string `env:"ADDRESS"`
}

func LoadConfig() *Config {
	var cfg Config
	return &cfg
}
