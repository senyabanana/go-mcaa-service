package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	flagRunAddr   string
	flagRunPoll   int
	flagRunReport int
)

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "HTTP server address")
	flag.IntVar(&flagRunPoll, "p", 2, "Poll interval in seconds")
	flag.IntVar(&flagRunReport, "r", 10, "Report interval in seconds")
	flag.Parse()

	// Проверка на наличие неизвестных флагов
	if flag.NFlag() == 0 && len(os.Args) > 1 {
		fmt.Println("Unknown flag provided")
		flag.Usage()
		os.Exit(1)
	}

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		flagRunAddr = envAddress
	}
	if envPoll := os.Getenv("POLL_INTERVAL"); envPoll != "" {
		value, err := strconv.Atoi(envPoll)
		if err != nil {
			log.Fatal(err)
		}
		flagRunPoll = value
	}
	if envReport := os.Getenv("REPORT_INTERVAL"); envReport != "" {
		value, err := strconv.Atoi(envReport)
		if err != nil {
			log.Fatal(err)
		}
		flagRunReport = value
	}
}
