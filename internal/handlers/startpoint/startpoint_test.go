package startpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandleStart(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		target  string
		method  string
		storage storage.Repository
		want    want
	}{
		{
			name:   "getAllMetrics",
			target: "/",
			method: http.MethodGet,
			storage: func() storage.Repository {
				s := storage.NewMemStorage(300, "", false)
				s.SetGauge("test_gauge", 1.2)
				s.SetCounter("test_counter", 10)
				return s
			}(),
			want: want{
				code:        http.StatusOK,
				contentType: "text/html",
			},
		},
		{
			name:    "noMetrics",
			target:  "/",
			method:  http.MethodGet,
			storage: storage.NewMemStorage(300, "", false),
			want: want{
				code:        http.StatusOK,
				contentType: "text/html",
			},
		},
		{
			name:    "invalidMethod",
			target:  "/",
			method:  http.MethodPost,
			storage: storage.NewMemStorage(300, "", false),
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый chi.Router
			r := chi.NewRouter()
			r.Get("/", HandleStart(tt.storage))

			// Создаем HTTP запрос
			req := httptest.NewRequest(tt.method, tt.target, nil)
			w := httptest.NewRecorder()

			// Обрабатываем запрос через chi.Router
			r.ServeHTTP(w, req)

			// Проверяем результаты
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
