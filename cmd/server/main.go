package main

import (
	"net/http"

	"github.com/senyabanana/go-mcaa-service/internal/agent/config"
	"github.com/senyabanana/go-mcaa-service/internal/server"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()
	config.ParseFlags(cfg)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	memStorage := storage.NewMemStorage()
	r := server.CreateServer(memStorage)
	logrus.Infof("Running server on %s\n", cfg.Address)
	logrus.Fatal(http.ListenAndServe(cfg.Address, r))
}
