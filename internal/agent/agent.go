package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/senyabanana/go-mcaa-service/internal/middleware"
	"github.com/senyabanana/go-mcaa-service/internal/models"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/sirupsen/logrus"
)

// Agent (HTTP-клиент) для сбора runtime-метрик и их последующей отправки на сервер.
type Agent struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	ServerURL      string
	gauges         map[string]float64
	counters       map[string]int64
	// мьютекс для синхронизации доступа к метрикам
	mu sync.Mutex
	// группа ожидания для создания горутин
	wg sync.WaitGroup
}

// NewAgent создает новый экземпляр агента.
func NewAgent(url string, pollInterval, reportInterval time.Duration) *Agent {
	return &Agent{
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		ServerURL:      url,
		gauges:         make(map[string]float64),
		counters:       make(map[string]int64),
	}
}

// RunAgent запускает агент и его горутины.
func (a *Agent) RunAgent() {
	// увеличиваем счетчик для горутин
	a.wg.Add(2)
	go a.collectMetrics()
	go a.sendMetrics()
	logrus.Info("Starting agent...")
	// ожидаем завершение всех горутин
	a.wg.Wait()
}

// collectMetric отвечает за сбор метрик С ЗАДАННОЙ ЧАСТОТОЙ.
func (a *Agent) collectMetrics() {
	// уменьшаем счетчик, когда метод завершится
	defer a.wg.Done()
	for {
		a.mu.Lock()
		a.collectRuntimeMetrics()
		logrus.Infof("%d ", a.counters["PollCount"])
		a.mu.Unlock()
		time.Sleep(a.PollInterval)
	}
}

// sendMetrics отвечает за отправку собранных метрик на сервер С ЗАДАННОЙ ЧАСТОТОЙ.
func (a *Agent) sendMetrics() {
	defer a.wg.Done()
	for {
		time.Sleep(a.ReportInterval)
		a.sendAllMetrics()
	}
}

// sendMetric отвечает за отправку одной метрики на сервер.
func (a *Agent) sendMetric(metric models.Metrics) error {
	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	b, err := middleware.Compress(data)
	if err != nil {
		return err
	}

	url := a.ServerURL + "/update"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", resp.Status)
	}

	// выполняем успешное логирование
	if metric.MType == storage.Gauge && metric.Value != nil {
		logrus.Infof("Successfully sent %s metric %s with value %v to %s", metric.MType, metric.ID, *metric.Value, url)
	} else if metric.MType == storage.Counter && metric.Delta != nil {
		logrus.Infof("Successfully sent %s metric %s with delta %v to %s", metric.MType, metric.ID, *metric.Delta, url)
	}

	return nil
}

// sendAllMetrics отвечает за отпраку всех собранных метрик на сервер БЕЗ ЗАДАННОЙ ЧАСТОТЫ.
func (a *Agent) sendAllMetrics() {
	a.mu.Lock()
	// создаем копии, чтобы избежать изменений в оригинальных данных во время отправки
	gaugesCopy := make(map[string]float64)
	countersCopy := make(map[string]int64)

	for name, value := range a.gauges {
		gaugesCopy[name] = value
	}
	for name, value := range a.counters {
		countersCopy[name] = value
	}
	logrus.Infof("===== Current pollcount is %d =====\n", a.counters["PollCount"])
	a.mu.Unlock()

	for name, value := range gaugesCopy {
		metric := models.Metrics{
			ID:    name,
			MType: storage.Gauge,
			Value: &value,
		}
		err := a.sendMetric(metric)
		if err != nil {
			logrus.Infof("failed to send %s: %v\n", name, err)
			return
		}
	}
	for name, value := range countersCopy {
		metric := models.Metrics{
			ID:    name,
			MType: storage.Counter,
			Delta: &value,
		}
		err := a.sendMetric(metric)
		if err != nil {
			logrus.Infof("failed to send %s: %v\n", name, err)
			return
		}
	}
}
