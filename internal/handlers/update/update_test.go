package update

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandleUpdate(t *testing.T) {
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
				code:        http.StatusMovedPermanently,
				contentType: "",
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
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый ServeMux
			mux := http.NewServeMux()
			mux.HandleFunc("/update/{type}/{name}/{value}", HandleUpdate(tt.storage))

			// Создаем HTTP запрос
			req := httptest.NewRequest(tt.method, tt.target, nil)
			w := httptest.NewRecorder()

			// Обрабатываем запрос через ServeHTTP
			mux.ServeHTTP(w, req)

			// Проверяем результаты
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
