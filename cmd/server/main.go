package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/senyabanana/go-mcaa-service/internal/handlers/startpoint"
	"github.com/senyabanana/go-mcaa-service/internal/handlers/update"
	"github.com/senyabanana/go-mcaa-service/internal/handlers/value"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
)

func main() {
	memStorage := storage.NewMemStorage()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route(`/`, func(r chi.Router) {
		r.Get(`/`, startpoint.HandleStart(memStorage))
		r.Route(`/update`, func(r chi.Router) {
			r.Post(`/{type}/{name}/{value}`, update.HandleUpdate(memStorage))
		})
		r.Route(`/value`, func(r chi.Router) {
			r.Get(`/{type}/{name}`, value.HandleValue(memStorage))
		})
	})
	log.Fatal(http.ListenAndServe(`:8080`, r))
}
