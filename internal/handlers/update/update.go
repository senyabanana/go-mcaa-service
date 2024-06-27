package update

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
)

// HandleUpdate обрабатывает HTTP-запроосы для обновления метрик.
func HandleUpdate(memStorage storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		rw.Header().Set("Content-Type", "text/plain")

		if r.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if metricName == "" {
			http.Error(rw, "metric name not provided", http.StatusNotFound)
			return
		}

		switch metricType {
		case storage.Gauge:
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(rw, "invalid value for gauge", http.StatusBadRequest)
				return
			}
			memStorage.UpdateGauge(metricName, value)
		case storage.Counter:
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(rw, "invalid value for counter", http.StatusBadRequest)
				return
			}
			memStorage.UpdateCounter(metricName, value)
		default:
			http.Error(rw, "unknown metric type", http.StatusBadRequest)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}
