package update

import (
	"encoding/json"
	"net/http"

	"github.com/senyabanana/go-mcaa-service/internal/models"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
)

// HandleUpdateJSON обрабатывает HTTP POST запросы на обновление метрик в формате JSON.
func HandleUpdateJSON(memStorage storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var m models.Metrics

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			http.Error(rw, "invalid request body", http.StatusBadRequest)
			return
		}

		switch m.MType {
		case storage.Gauge:
			if m.Value == nil {
				http.Error(rw, "missing value for gauge", http.StatusBadRequest)
				return
			}
			memStorage.UpdateGauge(m.ID, *m.Value)
		case storage.Counter:
			if m.Delta == nil {
				http.Error(rw, "missing value for counter", http.StatusBadRequest)
				return
			}
			memStorage.UpdateCounter(m.ID, *m.Delta)
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
