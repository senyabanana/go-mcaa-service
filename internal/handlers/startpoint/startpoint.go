package startpoint

import (
	"net/http"

	"github.com/senyabanana/go-mcaa-service/internal/storage"
)

// HandleStart выводит список метрик с их значением, известных на данный момент.
func HandleStart(memStorage storage.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html")

		if r.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		metrics := memStorage.GetAllMetrics()
		rw.Write([]byte(metrics))
		rw.WriteHeader(http.StatusOK)
	}
}
