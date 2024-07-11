package storage

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

// AllMetrics хранит все метрики для сериализации
type AllMetrics struct {
	Gauge   map[string]float64 `json:"gauges"`
	Counter map[string]int64   `json:"counters"`
}

// saveStorageToFile сохраняет текущие метрики в файл
func saveStorageToFile(s *MemStorage, filePath string) error {
	metrics := AllMetrics{
		Counter: s.counters,
		Gauge:   s.gauges,
	}
	// Проверка и создание файла, если его нет
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

// Dump периодически сохраняет метрики из хранилища в файл.
func Dump(s *MemStorage, filePath string, storeInterval int) error {
	dir, _ := path.Split(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logrus.Info(err)
		}
	}

	// Проверка и создание файла, если его нет
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err := os.Create(filePath)
		if err != nil {
			logrus.Info(err)
		}
	}

	ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		err := saveStorageToFile(s, filePath)
		if err != nil {
			logrus.Info(err)
		}
	}
	return nil
}

// LoadStorageFromFile загружает метрики из файла в хранилище
func LoadStorageFromFile(s *MemStorage, filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Info(err)
		return err
	}

	var data AllMetrics
	if err := json.Unmarshal(file, &data); err != nil {
		logrus.Info(err)
		return err
	}

	for name, value := range data.Counter {
		s.UpdateCounter(name, value)
	}
	for name, value := range data.Gauge {
		s.UpdateGauge(name, value)
	}

	return nil
}
