package agent

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/senyabanana/go-mcaa-service/internal/storage"
)

// Agent (HTTP-клиент) для сбора runtime-метрик и их последующей отправки на сервер
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

// NewAgent создает новый экземпляр агента
func NewAgent(url string, pollInterval, reportInterval time.Duration) *Agent {
	return &Agent{
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		ServerUrl:      url,
		gauges:         make(map[string]float64),
		counters:       make(map[string]int64),
	}
}

// RunAgent запускает агент и его горутины
func (a *Agent) RunAgent() {
	// увеличиваем счетчик для горутин
	a.wg.Add(2)
	go a.collectMetrics()
	go a.sendMetrics()
	fmt.Println("Starting agent...")
	// ожидаем завершение всех горутин
	a.wg.Wait()
}

// collectMetric отвечает за сбор метрик С ЗАДАННОЙ ЧАСТОТОЙ
func (a *Agent) collectMetrics() {
	// уменьшаем счетчик, когда метод завершится
	defer a.wg.Done()
	for {
		a.mu.Lock()
		a.collectRuntimeMetrics()
		fmt.Printf("%d ", a.counters["PollCount"])
		a.mu.Unlock()
		time.Sleep(a.PollInterval)
	}
}

// sendMetrics отвечает за отправку собранных метрик на сервер С ЗАДАННОЙ ЧАСТОТОЙ
func (a *Agent) sendMetrics() {
	defer a.wg.Done()
	for {
		time.Sleep(a.ReportInterval)
		a.sendAllMetrics()
	}
}

// collectRuntimeMetrics собирает runtime-метрики и записывает их в мапу
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

// sendOnceMetric отвечает за отправку одной метрики на сервер
func (a *Agent) sendMetric(metricType, name string, value interface{}) error {
	url := fmt.Sprintf("%s/update/%s/%s/%v", a.ServerUrl, metricType, name, value)
	log.Printf("sending metric: %s/%s to %s with value: %v\n", metricType, name, a.ServerUrl, value)
	response, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return err
	}
	err = response.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

// sendAllMetrics отвечает за отпраку всех собранных метрик на сервер БЕЗ ЗАДАННОЙ ЧАСТОТЫ
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
	fmt.Printf("\n===== Current pollcount is %d =====\n", a.counters["PollCount"])
	a.mu.Unlock()

	for name, value := range gaugesCopy {
		err := a.sendMetric(storage.Gauge, name, value)
		if err != nil {
			log.Printf("failed to send %s: %v\n", name, err)
			return
		}
	}
	for name, value := range countersCopy {
		err := a.sendMetric(storage.Counter, name, value)
		if err != nil {
			log.Printf("failed to send %s: %v\n", name, err)
			return
		}
	}
}
