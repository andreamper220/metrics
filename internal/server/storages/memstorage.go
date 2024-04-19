package storages

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/andreamper220/metrics.git/internal/shared"
)

type MemStorage struct {
	Counters map[shared.CounterMetricName]int64
	Gauges   map[shared.GaugeMetricName]float64
}

func (ms *MemStorage) WriteMetrics() {
	var mstats runtime.MemStats
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

func sendMetric(url, mType, name, value string, client *http.Client) error {
	requestURL := fmt.Sprintf("%s/update/%s/%s/%s", url, mType, name, value)

	res, err := client.Post(requestURL, "text/plain", http.NoBody)
	if err != nil {
		return err
	}

	return res.Body.Close()
}

func (ms *MemStorage) SendMetrics(url string) error {
	var err error = nil
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for name, value := range ms.Gauges {
		err = sendMetric(url, shared.GaugeMetricType, string(name), fmt.Sprintf("%f", value), client)
	}
	for name, value := range ms.Counters {
		err = sendMetric(url, shared.CounterMetricType, string(name), fmt.Sprintf("%d", value), client)
	}

	// return last error
	return err
}