package main

import (
	storage "github.com/senyabanana/go-mcaa-service/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	memStorage := storage.NewMemStorage()
	http.HandleFunc(`/update/`, func(w http.ResponseWriter, r *http.Request) {
		handleUpdate(memStorage, w, r)
	})

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}

// handleUpdate обрабатывает HTTP-запроосы для обновления метрик
func handleUpdate(memStorage *storage.MemStorage, w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(urlParts) != 4 {
		http.Error(w, "invalid request", http.StatusNotFound)
		return
	}
	metricType, metricName, metricValue := urlParts[1], urlParts[2], urlParts[3]
	if metricName == "" {
		http.Error(w, "missing metric name", http.StatusNotFound)
		return
	}
	switch metricType {
	case storage.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "invalid gauge value", http.StatusBadRequest)
			return
		}
		memStorage.UpdateGauge(metricValue, value)
	case storage.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
			return
		}
		memStorage.UpdateCounter(metricValue, value)
	default:
		http.Error(w, "invalid metric type", http.StatusBadRequest)
	}
	responseBody := "OK"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseBody))
}
