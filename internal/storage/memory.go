package storage

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

// saveStorageToFile сохраняет текущие метрики в файл в формате JSON.
func saveStorageToFile(ms *MemStorage, filePath string) error {
	var metrics AllMetrics
	metrics.Counter = ms.counters
	metrics.Gauge = ms.gauges

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0666)

}

// Dump периодически сохраняет метрики в файл с указанным интервалом.
func Dump(ms *MemStorage, filePath string, storeInterval int) error {
	dir, _ := path.Split(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			logrus.Info(err)
		}
	}
	pollTicker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer pollTicker.Stop()
	for range pollTicker.C {
		err := saveStorageToFile(ms, filePath)
		if err != nil {
			logrus.Info(err)
		}
	}
	return nil
}

// LoadStorageFromFile загружает метрики из файла и обновляет текущее состояние метрик в памяти.
func LoadStorageFromFile(ms *MemStorage, filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Info(err)
	}

	var data AllMetrics
	if err := json.Unmarshal(file, &data); err != nil {
		logrus.Info(err)
	}

	if len(data.Counter) != 0 {
		ms.UpdateCounterData(data.Counter)
	}
	if len(data.Gauge) != 0 {
		ms.UpdateGaugeData(data.Gauge)
	}
	return err
}
