package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData  *responseData
	headerWritten bool // проверяем, был ли заголовок уже записан
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	if !r.headerWritten {
		r.WriteHeader(http.StatusOK)
	}
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	if !r.headerWritten {
		r.ResponseWriter.WriteHeader(statusCode)
		r.responseData.status = statusCode
		r.headerWritten = true
	}
}

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
