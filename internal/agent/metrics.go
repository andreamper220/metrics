package agent

import (
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/andreamper220/metrics.git/internal/shared"
)

type metricsStruct struct {
	Counters map[shared.CounterMetricName]int64
	Gauges   map[shared.GaugeMetricName]float64
}

var metrics = &metricsStruct{
	Counters: make(map[shared.CounterMetricName]int64),
	Gauges:   make(map[shared.GaugeMetricName]float64),
}

func updateMetrics() {
	var mstats runtime.MemStats

	pollTicker := time.NewTicker(time.Duration(Config.PollInterval) * time.Second)
	for range pollTicker.C {
		runtime.ReadMemStats(&mstats)

		metrics.Gauges[shared.Alloc] = float64(mstats.Alloc)
		metrics.Gauges[shared.BuckHashSys] = float64(mstats.BuckHashSys)
		metrics.Gauges[shared.Frees] = float64(mstats.Frees)
		metrics.Gauges[shared.GcCPUFraction] = mstats.GCCPUFraction
		metrics.Gauges[shared.GcSys] = float64(mstats.GCSys)
		metrics.Gauges[shared.HeapAlloc] = float64(mstats.HeapAlloc)
		metrics.Gauges[shared.HeapIdle] = float64(mstats.HeapIdle)
		metrics.Gauges[shared.HeapInuse] = float64(mstats.HeapInuse)
		metrics.Gauges[shared.HeapObjects] = float64(mstats.HeapObjects)
		metrics.Gauges[shared.HeapReleased] = float64(mstats.HeapReleased)
		metrics.Gauges[shared.HeapSys] = float64(mstats.HeapSys)
		metrics.Gauges[shared.LastGc] = float64(mstats.LastGC)
		metrics.Gauges[shared.Lookups] = float64(mstats.Lookups)
		metrics.Gauges[shared.MemCacheInuse] = float64(mstats.MCacheInuse)
		metrics.Gauges[shared.MemCacheSys] = float64(mstats.MCacheSys)
		metrics.Gauges[shared.MemSpanInuse] = float64(mstats.MSpanInuse)
		metrics.Gauges[shared.MemSpanSys] = float64(mstats.MSpanSys)
		metrics.Gauges[shared.MemAllocs] = float64(mstats.Mallocs)
		metrics.Gauges[shared.NextGc] = float64(mstats.NextGC)
		metrics.Gauges[shared.NumForcedGc] = float64(mstats.NumForcedGC)
		metrics.Gauges[shared.NumGc] = float64(mstats.NumGC)
		metrics.Gauges[shared.OtherSys] = float64(mstats.OtherSys)
		metrics.Gauges[shared.PauseTotalNs] = float64(mstats.PauseTotalNs)
		metrics.Gauges[shared.RandomValue] = rand.Float64()
		metrics.Gauges[shared.StackInuse] = float64(mstats.StackInuse)
		metrics.Gauges[shared.StackSys] = float64(mstats.StackSys)
		metrics.Gauges[shared.Sys] = float64(mstats.Sys)
		metrics.Gauges[shared.TotalAlloc] = float64(mstats.TotalAlloc)

		metrics.Counters[shared.PollCount] = 1
	}
}

func updatePsUtilsMetrics() {
	pollTicker := time.NewTicker(time.Duration(Config.PollInterval) * time.Second)
	for range pollTicker.C {
		memory, _ := mem.VirtualMemory()
		cpuUtilization, _ := cpu.Percent(0, false)

		metrics.Gauges[shared.TotalMemory] = float64(memory.Total)
		metrics.Gauges[shared.FreeMemory] = float64(memory.Free)
		for i, cpuUtil := range cpuUtilization {
			metrics.Gauges[shared.GaugeMetricName("CPUutilization"+strconv.Itoa(i))] = cpuUtil
		}
	}
}

func readMetrics() *metricsStruct {
	return metrics
}
