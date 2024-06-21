package main

import (
	"net/http"

	"github.com/senyabanana/go-mcaa-service/internal/handlers/update"
	"github.com/senyabanana/go-mcaa-service/internal/storage/memory"
)

func main() {
	memStorage := memory.NewMemStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{type}/{name}/{value}", update.HandleUpdate(memStorage))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
