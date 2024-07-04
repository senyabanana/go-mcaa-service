package update

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
)

// HandleUpdatePlain обрабатывает HTTP-запроосы для обновления метрик.
func HandleUpdatePlain(memStorage storage.Repository) http.HandlerFunc {
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

//
//// HandleUpdateJSON передает метрики методом POST.
//func HandleUpdateJSON(memStorage storage.Repository) http.HandlerFunc {
//	return func(rw http.ResponseWriter, r *http.Request) {
//		rw.Header().Set("Content-Type", "application/json")
//
//		var m storage.Metrics
//		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
//			http.Error(rw, "invalid request body", http.StatusBadRequest)
//			return
//		}
//
//		switch m.MType {
//		case storage.Gauge:
//			if m.Value == nil {
//				http.Error(rw, "missing value for gauge", http.StatusBadRequest)
//				return
//			}
//			memStorage.UpdateGauge(m.ID, *m.Value)
//		case storage.Counter:
//			if m.Delta == nil {
//				http.Error(rw, "missing delta for counter", http.StatusBadRequest)
//				return
//			}
//			memStorage.UpdateCounter(m.ID, *m.Delta)
//		default:
//			http.Error(rw, "unknown metric type", http.StatusBadRequest)
//			return
//		}
//
//		resp, err := json.Marshal(m)
//		if err != nil {
//			http.Error(rw, "could not marshal response", http.StatusInternalServerError)
//			return
//		}
//
//		rw.WriteHeader(http.StatusOK)
//		rw.Write(resp)
//	}
//}
