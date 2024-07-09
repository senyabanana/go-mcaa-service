package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// responseData содержит информацию о статусе ответа и его размере.
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter оборачивает стандартный http.ResponseWriter,
// добавляя возможность логирования статуса и размера ответа.
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData  *responseData
	headerWritten bool // проверяем, был ли заголовок уже записан
}

// Write переопределяет метод Write, чтобы записывать размер данных в responseData.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	if !r.headerWritten {
		r.WriteHeader(http.StatusOK)
	}
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader переопределяет метод WriteHeader, чтобы записывать статус ответа в responseData.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	if !r.headerWritten {
		r.ResponseWriter.WriteHeader(statusCode)
		r.responseData.status = statusCode
		r.headerWritten = true
	}
}

// LoggingMiddleware оборачивает следующий http.Handler,
// логируя информацию о каждом запросе.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: rw,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r)
		duration := time.Since(start)
		logrus.WithFields(logrus.Fields{
			"uri":      r.RequestURI,
			"method":   r.Method,
			"status":   responseData.status,
			"duration": duration,
			"size":     responseData.size,
		}).Info("Handled request")
	})
}
