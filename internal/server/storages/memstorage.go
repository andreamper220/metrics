package storages

import (
	"encoding/json"
	"math/rand"
	"os"
	"runtime"

	"github.com/andreamper220/metrics.git/internal/shared"
)

var FileStoragePath string
var ToSaveMetricsAsync bool = false

type MemStorage struct {
	Counters map[shared.CounterMetricName]int64
	Gauges   map[shared.GaugeMetricName]float64
}

var Storage = &MemStorage{
	Counters: make(map[shared.CounterMetricName]int64),
	Gauges:   make(map[shared.GaugeMetricName]float64),
}

func (ms *MemStorage) WriteMetrics() {
	var mstats runtime.MemStats
	runtime.ReadMemStats(&mstats)

	ms.Gauges = map[shared.GaugeMetricName]float64{
		shared.Alloc:         float64(mstats.Alloc),
		shared.BuckHashSys:   float64(mstats.BuckHashSys),
		shared.Frees:         float64(mstats.Frees),
		shared.GcCPUFraction: mstats.GCCPUFraction,
		shared.GcSys:         float64(mstats.GCSys),
		shared.HeapAlloc:     float64(mstats.HeapAlloc),
		shared.HeapIdle:      float64(mstats.HeapIdle),
		shared.HeapInuse:     float64(mstats.HeapInuse),
		shared.HeapObjects:   float64(mstats.HeapObjects),
		shared.HeapReleased:  float64(mstats.HeapReleased),
		shared.HeapSys:       float64(mstats.HeapSys),
		shared.LastGc:        float64(mstats.LastGC),
		shared.Lookups:       float64(mstats.Lookups),
		shared.MemCacheInuse: float64(mstats.MCacheInuse),
		shared.MemCacheSys:   float64(mstats.MCacheSys),
		shared.MemSpanInuse:  float64(mstats.MSpanInuse),
		shared.MemSpanSys:    float64(mstats.MSpanSys),
		shared.MemAllocs:     float64(mstats.Mallocs),
		shared.NextGc:        float64(mstats.NextGC),
		shared.NumForcedGc:   float64(mstats.NumForcedGC),
		shared.NumGc:         float64(mstats.NumGC),
		shared.OtherSys:      float64(mstats.OtherSys),
		shared.PauseTotalNs:  float64(mstats.PauseTotalNs),
		shared.RandomValue:   rand.Float64(),
		shared.StackInuse:    float64(mstats.StackInuse),
		shared.StackSys:      float64(mstats.StackSys),
		shared.Sys:           float64(mstats.Sys),
		shared.TotalAlloc:    float64(mstats.TotalAlloc),
	}
	ms.Counters[shared.PollCount] = 1
}

func (ms *MemStorage) StoreMetrics() error {
	file, err := os.OpenFile(FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	data := make([]byte, 0)
	for name, value := range Storage.Counters {
		metric := &shared.Metric{
			ID:    string(name),
			MType: shared.CounterMetricType,
			Delta: &value,
		}
		metricData, _ := json.Marshal(&metric)
		metricData = append(metricData, '\n')
		data = append(data, metricData...)
	}

	for name, value := range Storage.Gauges {
		metric := &shared.Metric{
			ID:    string(name),
			MType: shared.GaugeMetricType,
			Value: &value,
		}
		metricData, _ := json.Marshal(&metric)
		metricData = append(metricData, '\n')
		data = append(data, metricData...)
	}

	_, err = file.Write(data)

	return err
}
