package main

import (
	"time"

	"github.com/senyabanana/go-mcaa-service/internal/agent"
)

func main() {
	parseFlags()

	url := "http://" + flagRunAddr
	pollInterval := time.Duration(flagRunPoll) * time.Second
	reportInterval := time.Duration(flagRunReport) * time.Second

	client := agent.NewAgent(url, pollInterval, reportInterval)
	client.RunAgent()
}
