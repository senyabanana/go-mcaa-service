package config

import (
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func ParseFlags(cfg *Config) {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "HTTP server address")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "interval for saving metrics on the server")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "path file storage to save data")
	flag.BoolVar(&cfg.Restore, "r", true, "need to load data at startup")
	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		cfg.Address = envAddress
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		value, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			logrus.Info(err)
		}
		cfg.StoreInterval = value
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		value, err := strconv.ParseBool(envRestore)
		if err != nil {
			logrus.Info(err)
		}
		cfg.Restore = value
	}
}
