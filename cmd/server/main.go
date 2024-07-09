package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/senyabanana/go-mcaa-service/internal/handlers/startpoint"
	"github.com/senyabanana/go-mcaa-service/internal/handlers/update"
	"github.com/senyabanana/go-mcaa-service/internal/handlers/value"
	"github.com/senyabanana/go-mcaa-service/internal/middleware"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	parseFlags()

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	memStorage := storage.NewMemStorage()

	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.GzipHandler)
	r.Route(`/`, func(r chi.Router) {
		r.Get(`/`, startpoint.HandleStart(memStorage))
		r.Route(`/update`, func(r chi.Router) {
			r.Post(`/{type}/{name}/{value}`, update.HandleUpdatePlain(memStorage))
			r.Post(`/`, update.HandleUpdateJSON(memStorage))
		})
		r.Route(`/value`, func(r chi.Router) {
			r.Get(`/{type}/{name}`, value.HandleValuePlain(memStorage))
			r.Post(`/`, value.HandleValueJSON(memStorage))
		})
	})
	logrus.Infof("Running server on %s\n", flagRunAddr)
	logrus.Fatal(http.ListenAndServe(flagRunAddr, r))
}
