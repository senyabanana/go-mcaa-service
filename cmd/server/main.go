package main

import (
	"net/http"

	handler "github.com/senyabanana/go-mcaa-service/internal/handler"
	storage "github.com/senyabanana/go-mcaa-service/internal/storage"
)

func main() {
	memStorage := storage.NewMemStorage()

	http.HandleFunc("POST /update/{type}/{name}/{value}", func(wr http.ResponseWriter, r *http.Request) {
		handler.HandleUpdate(memStorage, wr, r)
	})

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
