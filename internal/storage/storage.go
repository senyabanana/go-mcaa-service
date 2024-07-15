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
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
}

// MemStorage - структура для хранения метрик.
type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

// NewMemStorage - создание нового экземпляра MemStorage.
func NewMemStorage(storeInterval int, filePath string, restore bool) *MemStorage {
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

// SetGauge устанавливает значение метрики типа gauge.
func (ms *MemStorage) SetGauge(name string, value float64) {
	ms.gauges[name] = value
}

// SetCounter устанавливает значение метрики типа counter.
func (ms *MemStorage) SetCounter(name string, value int64) {
	ms.counters[name] = value
}

// AllMetrics содержит все метрики.
type AllMetrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// AllMetrics позволяет получить все метрики в одном объекте.
func (ms *MemStorage) AllMetrics() *AllMetrics {
	return &AllMetrics{
		Gauge:   ms.gauges,
		Counter: ms.counters,
	}
}

// UpdateGaugeData позволяет обновлять данные для всех ключей.
func (ms *MemStorage) UpdateGaugeData(gaugeData map[string]float64) {
	ms.gauges = gaugeData
}

// UpdateCounterData позволяет обновлять данные для всех ключей.
func (ms *MemStorage) UpdateCounterData(counterData map[string]int64) {
	ms.counters = counterData
}
