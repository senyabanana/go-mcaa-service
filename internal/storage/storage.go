package storage

// Repository имплементирует интерфейс хранения
type Repository interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
}
