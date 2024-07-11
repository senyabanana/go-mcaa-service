package main

import (
	"time"

	"github.com/senyabanana/go-mcaa-service/internal/agent"
	"github.com/senyabanana/go-mcaa-service/internal/agent/config"
)

func main() {
	//parseFlags()
	cfg := config.LoadConfig()
	config.ParseFlags(cfg)

	url := "http://" + cfg.Address
	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	client := agent.NewAgent(url, pollInterval, reportInterval)
	client.RunAgent()
}
