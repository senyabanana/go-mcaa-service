package agent

import (
	"math/rand"
	"runtime"
)

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
