package main

import (
	"net/http"

	handler "github.com/senyabanana/go-mcaa-service/internal/handler"
	storage "github.com/senyabanana/go-mcaa-service/internal/storage"
)

func main() {
	memStorage := storage.NewMemStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{type}/{name}/{value}", func(rw http.ResponseWriter, r *http.Request) {
		handler.HandleUpdate(memStorage, rw, r)
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
