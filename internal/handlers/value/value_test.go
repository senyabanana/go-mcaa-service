package value

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/senyabanana/go-mcaa-service/internal/models"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandleValuePlain(t *testing.T) {
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
				s := storage.NewMemStorage(300, "", false)
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
				s := storage.NewMemStorage(300, "", false)
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
			storage: storage.NewMemStorage(300, "", false),
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
			storage: storage.NewMemStorage(300, "", false),
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
			storage: storage.NewMemStorage(300, "", false),
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
			r.Get("/value/{type}/{name}", HandleValuePlain(tt.storage))

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

func TestHandleValueJSON(t *testing.T) {
	memStorage := storage.NewMemStorage(300, "", false)

	// Предварительное заполнение хранилища метриками
	initialGauge := 123.45
	initialCounter := int64(123)
	memStorage.UpdateGauge("TestGauge", initialGauge)
	memStorage.UpdateCounter("TestCounter", initialCounter)

	tests := []struct {
		name           string
		requestMetric  models.Metrics
		expectedMetric models.Metrics
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid Gauge Metric",
			requestMetric: models.Metrics{
				ID:    "TestGauge",
				MType: "gauge",
			},
			expectedMetric: models.Metrics{
				ID:    "TestGauge",
				MType: "gauge",
				Value: &initialGauge,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid Counter Metric",
			requestMetric: models.Metrics{
				ID:    "TestCounter",
				MType: "counter",
			},
			expectedMetric: models.Metrics{
				ID:    "TestCounter",
				MType: "counter",
				Delta: &initialCounter,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Metric Not Found",
			requestMetric: models.Metrics{
				ID:    "NonExistent",
				MType: "gauge",
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "metric not found\n",
		},
		{
			name: "Unknown Metric Type",
			requestMetric: models.Metrics{
				ID:    "TestUnknown",
				MType: "unknown",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "unknown metric type\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := HandleValueJSON(memStorage)

			metricData, err := json.Marshal(tt.requestMetric)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/value", bytes.NewBuffer(metricData))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				var responseMetric models.Metrics
				err := json.NewDecoder(rr.Body).Decode(&responseMetric)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMetric, responseMetric)
			} else {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}
