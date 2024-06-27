package main

import (
	"time"

	"github.com/senyabanana/go-mcaa-service/internal/agent"
)

const (
	url            = "http://localhost:8080"
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	client := agent.NewAgent(url, pollInterval, reportInterval)
	client.RunAgent()
}
