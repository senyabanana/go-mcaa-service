package main

import (
	"net/http"

	"github.com/senyabanana/go-mcaa-service/internal/server"
	"github.com/senyabanana/go-mcaa-service/internal/server/config"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()
	config.ParseFlags(cfg)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	memStorage := storage.NewMemStorage(cfg.StoreInterval, cfg.FileStoragePath, cfg.Restore)
	if cfg.FileStoragePath != "" {
		if cfg.Restore {
			err := storage.LoadStorageFromFile(memStorage, cfg.FileStoragePath)
			if err != nil {
				logrus.Info(err)
			}
		}
		if cfg.StoreInterval != 0 {
			go func() {
				err := storage.Dump(memStorage, cfg.FileStoragePath, cfg.StoreInterval)
				if err != nil {
					logrus.Info(err)
				}
			}()
		}
	}
	r := server.CreateServer(memStorage)
	logrus.Infof("Running server on %s\n", cfg.Address)
	logrus.Fatal(http.ListenAndServe(cfg.Address, r))
}
