package update

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/senyabanana/go-mcaa-service/internal/models"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdatePlain(t *testing.T) {
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
			name:    "updateGauge",
			target:  "/update/gauge/gauge_test/1.2",
			method:  http.MethodPost,
			storage: storage.NewMemStorage(),
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "updateCounter",
			target:  "/update/counter/counter_test/10",
			method:  http.MethodPost,
			storage: storage.NewMemStorage(),
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "movedPermanently",
			target:  "/update/gauge//22.1",
			method:  http.MethodPost,
			storage: storage.NewMemStorage(),
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "invalidGauge",
			target:  "/update/gauge/my_gauge/sss",
			method:  http.MethodPost,
			storage: storage.NewMemStorage(),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "invalidCounter",
			target:  "/update/counter/my_counter/22.3",
			method:  http.MethodPost,
			storage: storage.NewMemStorage(),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "invalidMetricType",
			target:  "/update/sss/my_sss/ewq",
			method:  http.MethodPost,
			storage: storage.NewMemStorage(),
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:    "invalidMethod",
			target:  "/update/gauge/gauge_test/1.2",
			method:  http.MethodGet,
			storage: storage.NewMemStorage(),
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
			r.Post("/update/{type}/{name}/{value}", HandleUpdatePlain(tt.storage))

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

func TestHandleUpdateJSON(t *testing.T) {
	memStorage := storage.NewMemStorage()

	tests := []struct {
		name           string
		metric         models.Metrics
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid Gauge Metric",
			metric: models.Metrics{
				ID:    "TestGauge",
				MType: "gauge",
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"TestGauge","type":"gauge","value":123.45}`,
		},
		{
			name: "Valid Counter Metric",
			metric: models.Metrics{
				ID:    "TestCounter",
				MType: "counter",
				Delta: func() *int64 { v := int64(123); return &v }(),
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"TestCounter","type":"counter","delta":123}`,
		},
		{
			name: "Missing Value for Gauge",
			metric: models.Metrics{
				ID:    "TestGauge",
				MType: "gauge",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing value for gauge\n",
		},
		{
			name: "Missing Value for Counter",
			metric: models.Metrics{
				ID:    "TestCounter",
				MType: "counter",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing value for counter\n",
		},
		{
			name: "Unknown Metric Type",
			metric: models.Metrics{
				ID:    "TestUnknown",
				MType: "unknown",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "unknown metric type\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := HandleUpdateJSON(memStorage)

			metricData, err := json.Marshal(tt.metric)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/update", bytes.NewBuffer(metricData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
