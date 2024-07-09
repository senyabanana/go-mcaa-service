package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// gzipWriter оборачивает http.ResponseWriter, добавляя поддержку сжатия gzip.
type gzipWriter struct {
	http.ResponseWriter
	w io.Writer
}

// Write переопределяет метод Write, чтобы записывать данные в gzip.Writer.
func (gw gzipWriter) Write(b []byte) (int, error) {
	return gw.w.Write(b)
}

// GzipHandler обрабатывает HTTP-запросы, поддерживая сжатие и разжатие gzip.
func GzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(rw, "Bad Request", http.StatusBadRequest)
				return
			}
			defer reader.Close()
			r.Body = io.NopCloser(reader)
		}
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gzWriter := gzip.NewWriter(rw)
			defer gzWriter.Close()
			rw.Header().Set("Content-Encoding", "gzip")
			rw = gzipWriter{ResponseWriter: rw, w: gzWriter}
		}

		next.ServeHTTP(rw, r)
	})
}

// Compress сжимает данные в формате gzip.
func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		logrus.Errorf("failed to create gzip writer: %v", err)
		return nil, err
	}
	_, err = w.Write(data)
	if err != nil {
		logrus.Errorf("failed to write data to gzip writer: %v", err)
		return nil, err
	}
	err = w.Close()
	if err != nil {
		logrus.Errorf("failed to close gzip writer: %v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}
