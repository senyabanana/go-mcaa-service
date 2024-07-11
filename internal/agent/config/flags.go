package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func ParseFlags(cfg *Config) {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "HTTP server address")
	flag.IntVar(&cfg.PollInterval, "p", 2, "Poll interval in seconds")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "Report interval in seconds")
	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		cfg.Address = envAddress
	}
	if envPoll := os.Getenv("POLL_INTERVAL"); envPoll != "" {
		value, err := strconv.Atoi(envPoll)
		if err != nil {
			logrus.Info(err)
		}
		cfg.PollInterval = value
	}
	if envReport := os.Getenv("REPORT_INTERVAL"); envReport != "" {
		value, err := strconv.Atoi(envReport)
		if err != nil {
			logrus.Info(err)
		}
		cfg.ReportInterval = value
	}
}
