package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/senyabanana/go-mcaa-service/internal/models"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/sirupsen/logrus"
)

// Agent (HTTP-клиент) для сбора runtime-метрик и их последующей отправки на сервер.
type Agent struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	ServerUrl      string
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
		ServerUrl:      url,
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

// collectRuntimeMetrics собирает runtime-метрики и записывает их в мапу.
func (a *Agent) collectRuntimeMetrics() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	a.gauges["Alloc"] = float64(ms.Alloc)
	a.gauges["BuckHashSys"] = float64(ms.BuckHashSys)
	a.gauges["Frees"] = float64(ms.Frees)
	a.gauges["GCCPUFraction"] = ms.GCCPUFraction
	a.gauges["GCSys"] = float64(ms.GCSys)
	a.gauges["HeapAlloc"] = float64(ms.HeapAlloc)
	a.gauges["HeapIdle"] = float64(ms.HeapIdle)
	a.gauges["HeapInuse"] = float64(ms.HeapInuse)
	a.gauges["HeapObjects"] = float64(ms.HeapObjects)
	a.gauges["HeapReleased"] = float64(ms.HeapReleased)
	a.gauges["HeapSys"] = float64(ms.HeapSys)
	a.gauges["LastGC"] = float64(ms.LastGC)
	a.gauges["Lookups"] = float64(ms.Lookups)
	a.gauges["MCacheInuse"] = float64(ms.MCacheInuse)
	a.gauges["MCacheSys"] = float64(ms.MCacheSys)
	a.gauges["MSpanInuse"] = float64(ms.MSpanInuse)
	a.gauges["MSpanSys"] = float64(ms.MSpanSys)
	a.gauges["Mallocs"] = float64(ms.Mallocs)
	a.gauges["NextGC"] = float64(ms.NextGC)
	a.gauges["NumForcedGC"] = float64(ms.NumForcedGC)
	a.gauges["NumGC"] = float64(ms.NumGC)
	a.gauges["OtherSys"] = float64(ms.OtherSys)
	a.gauges["PauseTotalNs"] = float64(ms.PauseTotalNs)
	a.gauges["StackInuse"] = float64(ms.StackInuse)
	a.gauges["StackSys"] = float64(ms.StackSys)
	a.gauges["Sys"] = float64(ms.Sys)
	a.gauges["TotalAlloc"] = float64(ms.TotalAlloc)
	a.gauges["RandomValue"] = rand.Float64()
	a.counters["PollCount"]++
}

// sendMetric отвечает за отправку одной метрики на сервер.
func (a *Agent) sendMetric(metric models.Metrics) error {
	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	url := a.ServerUrl + "/update"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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
