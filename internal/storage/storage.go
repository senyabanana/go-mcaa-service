package storage

const (
	Gauge   = "gauge"
	Counter = "counter"
)

// MemStorage - структура для хранения метрик
type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

// NewMemStorage - создание нового экземпляра MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// UpdateGauge обновляет метрику типа gauges
func (ms *MemStorage) UpdateGauge(name string, value float64) {
	ms.gauges[name] = value
}

// UpdateCounter обновляет метрику типа counters
func (ms *MemStorage) UpdateCounter(name string, value int64) {
	ms.counters[name] += value
}
