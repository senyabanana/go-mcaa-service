package storage

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

func saveStorageToFile(s *MemStorage, filePath string) error {
	var metrics AllMetrics
	metrics.Counter = s.counters
	metrics.Gauge = s.gauges

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0666)

}

func Dump(s *MemStorage, filePath string, storeInterval int) error {
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
		err := saveStorageToFile(s, filePath)
		if err != nil {
			logrus.Info(err)
		}
	}
	return nil
}

func LoadStorageFromFile(s *MemStorage, filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Info(err)
	}

	var data AllMetrics
	if err := json.Unmarshal(file, &data); err != nil {
		logrus.Info(err)
	}

	if len(data.Counter) != 0 {
		s.UpdateCounterData(data.Counter)
	}
	if len(data.Gauge) != 0 {
		s.UpdateGaugeData(data.Gauge)
	}
	return err
}
