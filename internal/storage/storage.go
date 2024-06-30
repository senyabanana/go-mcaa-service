package storage

import (
	"fmt"
	"strings"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

// Repository имплементирует интерфейс хранения.
type Repository interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllMetrics() string
}

// MemStorage - структура для хранения метрик.
type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

// NewMemStorage - создание нового экземпляра MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// UpdateGauge обновляет метрику типа gauges.
func (ms *MemStorage) UpdateGauge(name string, value float64) {
	ms.gauges[name] = value
}

// UpdateCounter обновляет метрику типа counters.
func (ms *MemStorage) UpdateCounter(name string, value int64) {
	ms.counters[name] += value
}

// GetGauge возвращает значение матрики типа gauge.
func (ms *MemStorage) GetGauge(name string) (float64, bool) {
	value, ok := ms.gauges[name]
	return value, ok
}

// GetCounter возвращает значение метрики типа counter.
func (ms *MemStorage) GetCounter(name string) (int64, bool) {
	value, ok := ms.counters[name]
	return value, ok
}

// GetAllMetrics возвращает все метрики в виде HTML-страницы.
func (ms *MemStorage) GetAllMetrics() string {
	var result strings.Builder
	result.WriteString("<html><body>")
	for name, value := range ms.gauges {
		result.WriteString(fmt.Sprintf("<br>gauge/%s: %v</br>", name, value))
	}
	for name, value := range ms.counters {
		result.WriteString(fmt.Sprintf("<br>counter/%s: %v</br>", name, value))
	}
	result.WriteString("</body></html>")
	return result.String()
}
