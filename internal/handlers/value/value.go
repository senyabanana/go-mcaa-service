package value

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
)

// HandleValue обрабатывает GET-запросы и выводит текущее значение метрики.
func HandleValue(memStorage storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")

		rw.Header().Set("Content-Type", "text/plain")

		if r.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if metricName == "" {
			http.Error(rw, "metric name not providen", http.StatusNotFound)
			return
		}

		switch metricType {
		case storage.Gauge:
			value, ok := memStorage.GetGauge(metricName)
			if !ok {
				http.Error(rw, "metric not found", http.StatusNotFound)
				return
			}
			rw.Write([]byte(fmt.Sprintf("%v", value)))
		case storage.Counter:
			value, ok := memStorage.GetCounter(metricName)
			if !ok {
				http.Error(rw, "metric not found", http.StatusNotFound)
				return
			}
			rw.Write([]byte(fmt.Sprintf("%v", value)))
		default:
			http.Error(rw, "unknown metric type", http.StatusBadRequest)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}
