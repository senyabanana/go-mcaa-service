package update

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/senyabanana/go-mcaa-service/internal/storage/memory"
)

// HandleUpdate обрабатывает HTTP-запроосы для обновления метрик
func HandleUpdate(memStorage storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := r.PathValue("type")
		metricName := r.PathValue("name")
		metricValue := r.PathValue("value")

		rw.Header().Set("Content-Type", "text/plain")

		if r.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// выводим информацию о запросе
		fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)

		if metricName == "" {
			http.Error(rw, "metric name not provided", http.StatusNotFound)
			return
		}

		switch metricType {
		case memory.Gauge:
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(rw, "invalid value for gauge", http.StatusBadRequest)
				return
			}
			memStorage.UpdateGauge(metricName, value)
		case memory.Counter:
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
