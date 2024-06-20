package main

import (
	"fmt"
	"runtime"
	"time"
)

func collectMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Printf("Alloc: %v KB\n", memStats.Alloc/1024)
	fmt.Printf("BuckHashSys: %v KB\n", memStats.BuckHashSys/1024)
	fmt.Printf("Frees: %v\n", memStats.Frees)
	fmt.Printf("GCCPUFraction: %v\n", memStats.GCCPUFraction)
	fmt.Printf("GCSys: %v KB\n", memStats.GCSys/1024)
	fmt.Printf("HeapAlloc: %v KB\n", memStats.HeapAlloc/1024)
	fmt.Printf("HeapIdle: %v KB\n", memStats.HeapIdle/1024)
	fmt.Printf("HeapInuse: %v KB\n", memStats.HeapInuse/1024)
	fmt.Printf("HeapObjects: %v\n", memStats.HeapObjects)
	fmt.Printf("HeapReleased: %v KB\n", memStats.HeapReleased/1024)
	fmt.Printf("HeapSys: %v KB\n", memStats.HeapSys/1024)
	fmt.Printf("LastGC: %v\n", time.Unix(0, int64(memStats.LastGC)))
	fmt.Printf("Lookups: %v\n", memStats.Lookups)
	fmt.Printf("MCacheInuse: %v KB\n", memStats.MCacheInuse/1024)
	fmt.Printf("MCacheSys: %v KB\n", memStats.MCacheSys/1024)
	fmt.Printf("MSpanInuse: %v KB\n", memStats.MSpanInuse/1024)
	fmt.Printf("MSpanSys: %v KB\n", memStats.MSpanSys/1024)
	fmt.Printf("Mallocs: %v\n", memStats.Mallocs)
	fmt.Printf("NextGC: %v KB\n", memStats.NextGC/1024)
	fmt.Printf("NumForcedGC: %v\n", memStats.NumForcedGC)
	fmt.Printf("NumGC: %v\n", memStats.NumGC)
	fmt.Printf("OtherSys: %v KB\n", memStats.OtherSys/1024)
	fmt.Printf("PauseTotalNs: %v\n", memStats.PauseTotalNs)
	fmt.Printf("StackInuse: %v KB\n", memStats.StackInuse/1024)
	fmt.Printf("StackSys: %v KB\n", memStats.StackSys/1024)
	fmt.Printf("Sys: %v KB\n", memStats.Sys/1024)
	fmt.Printf("TotalAlloc: %v KB\n", memStats.TotalAlloc/1024)
}

func main() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			collectMetrics()
		}
	}
}
