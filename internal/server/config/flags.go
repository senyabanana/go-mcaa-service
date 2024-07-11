package config

import (
	"flag"
	"os"
)

func ParseFlags(cfg *Config) {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "HTTP server address")
	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		cfg.Address = envAddress
	}
}
