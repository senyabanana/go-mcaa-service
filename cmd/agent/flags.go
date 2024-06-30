package main

import (
	"flag"
	"fmt"
	"os"
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
}
