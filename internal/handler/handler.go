package handler

import (
	"net/http"
	"strconv"

	"github.com/senyabanana/go-mcaa-service/internal/storage"
)

// HandleUpdate обрабатывает HTTP-запроосы для обновления метрик
func HandleUpdate(memStorage *storage.MemStorage, wr http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("type")
	metricName := r.PathValue("name")
	metricValue := r.PathValue("value")

	if metricName == "" {
		http.Error(wr, "metric name not provided", http.StatusNotFound)
		return
	}

	switch metricType {
	case storage.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(wr, "invalid value for gauge", http.StatusBadRequest)
			return
		}
		memStorage.UpdateGauge(metricName, value)
	case storage.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(wr, "invalid value for counter", http.StatusBadRequest)
			return
		}
		memStorage.UpdateCounter(metricName, value)
	default:
		http.Error(wr, "unknown metric type", http.StatusBadRequest)
		return
	}

	wr.WriteHeader(http.StatusOK)
}
