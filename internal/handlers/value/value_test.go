package value

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandleValue(t *testing.T) {
	type want struct {
		code        int
		contentType string
		body        string
	}
	tests := []struct {
		name    string
		target  string
		method  string
		storage storage.Repository
		want    want
	}{
		{
			name:   "getGauge",
			target: "/value/gauge/test_gauge",
			method: http.MethodGet,
			storage: func() storage.Repository {
				s := storage.NewMemStorage()
				s.SetGauge("test_gauge", 1.2)
				return s
			}(),
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
				body:        "1.2",
			},
		},
		{
			name:   "getCounter",
			target: "/value/counter/test_counter",
			method: http.MethodGet,
			storage: func() storage.Repository {
				s := storage.NewMemStorage()
				s.SetCounter("test_counter", 10)
				return s
			}(),
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
				body:        "10",
			},
		},
		{
			name:    "metricNotFound",
			target:  "/value/gauge/unknown",
			method:  http.MethodGet,
			storage: storage.NewMemStorage(),
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				body:        "metric not found\n",
			},
		},
		{
			name:    "unknownMetricType",
			target:  "/value/unknown/test",
			method:  http.MethodGet,
			storage: storage.NewMemStorage(),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "unknown metric type\n",
			},
		},
		{
			name:    "invalidMethod",
			target:  "/value/gauge/test_gauge",
			method:  http.MethodPost,
			storage: storage.NewMemStorage(),
			want: want{
				code:        http.StatusMethodNotAllowed,
				contentType: "",
				body:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый chi.Router
			r := chi.NewRouter()
			r.Get("/value/{type}/{name}", HandleValue(tt.storage))

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

			bodyBytes, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.want.body, string(bodyBytes))
		})
	}
}
