package value

import (
	"encoding/json"
	"net/http"

	"github.com/senyabanana/go-mcaa-service/internal/models"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
)

// HandleValueJSON обрабатывает HTTP POST запросы на получение значений метрик в формате JSON.
func HandleValueJSON(memStorage storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var m models.Metrics

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(rw, "invalid request body", http.StatusBadRequest)
			return
		}

		switch m.MType {
		case storage.Gauge:
			value, ok := memStorage.GetGauge(m.ID)
			if !ok {
				http.Error(rw, "metric not found", http.StatusNotFound)
				return
			}
			m.Value = &value
		case storage.Counter:
			delta, ok := memStorage.GetCounter(m.ID)
			if !ok {
				http.Error(rw, "metric not found", http.StatusNotFound)
				return
			}
			m.Delta = &delta
		default:
			http.Error(rw, "unknown metric type", http.StatusBadRequest)
			return
		}

		resp, err := json.Marshal(m)
		if err != nil {
			http.Error(rw, "could not marshal response", http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(resp)
	}
}
