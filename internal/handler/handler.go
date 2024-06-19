package handler

import (
	"net/http"
	"strconv"
	"strings"

	storage "github.com/senyabanana/go-mcaa-service/internal/storage"
)

// HandleUpdate обрабатывает HTTP-запроосы для обновления метрик
func HandleUpdate(memStorage *storage.MemStorage, w http.ResponseWriter, r *http.Request) {
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
