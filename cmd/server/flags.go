package main

import (
	"flag"
	"fmt"
	"os"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "HTTP server address")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	// Проверка на наличие неизвестных флагов
	if flag.NFlag() == 0 && len(os.Args) > 1 {
		fmt.Println("Unknown flag provided")
		flag.Usage()
		os.Exit(1)
	}
}
