package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	_const "github.com/andreamper220/metrics.git/internal/const"
)

type MemStorage struct {
	counters map[_const.CounterMetricName]int64
	gauges   map[_const.GaugeMetricName]float64
}

func (ms *MemStorage) writeMetrics() {
	var mstats runtime.MemStats
	ms.gauges = map[_const.GaugeMetricName]float64{
		_const.Alloc:         float64(mstats.Alloc),
		_const.BuckHashSys:   float64(mstats.BuckHashSys),
		_const.Frees:         float64(mstats.Frees),
		_const.GcCpuFraction: mstats.GCCPUFraction,
		_const.GcSys:         float64(mstats.GCSys),
		_const.HeapAlloc:     float64(mstats.HeapAlloc),
		_const.HeapIdle:      float64(mstats.HeapIdle),
		_const.HeapInuse:     float64(mstats.HeapInuse),
		_const.HeapObjects:   float64(mstats.HeapObjects),
		_const.HeapReleased:  float64(mstats.HeapReleased),
		_const.HeapSys:       float64(mstats.HeapSys),
		_const.LastGc:        float64(mstats.LastGC),
		_const.Lookups:       float64(mstats.Lookups),
		_const.MemCacheInuse: float64(mstats.MCacheInuse),
		_const.MemCacheSys:   float64(mstats.MCacheSys),
		_const.MemSpanInuse:  float64(mstats.MSpanInuse),
		_const.MemSpanSys:    float64(mstats.MSpanSys),
		_const.MemAllocs:     float64(mstats.Mallocs),
		_const.NextGc:        float64(mstats.NextGC),
		_const.NumForcedGc:   float64(mstats.NumForcedGC),
		_const.NumGc:         float64(mstats.NumGC),
		_const.OtherSys:      float64(mstats.OtherSys),
		_const.PauseTotalNs:  float64(mstats.PauseTotalNs),
		_const.RandomValue:   rand.Float64(),
		_const.StackInuse:    float64(mstats.StackInuse),
		_const.StackSys:      float64(mstats.StackSys),
		_const.Sys:           float64(mstats.Sys),
		_const.TotalAlloc:    float64(mstats.TotalAlloc),
	}
	ms.counters[_const.PollCount] = 1
}

func sendMetric(url, name, value string) error {
	requestURL := fmt.Sprintf("%s/update/%s/%s/%s", url, _const.GaugeMetricType, name, value)
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	_, err = client.Do(req)
	return err
}

func (ms *MemStorage) sendMetrics(url string) error {
	var err error = nil

	for name, value := range ms.gauges {
		err = sendMetric(url, string(name), fmt.Sprintf("%f", value))
	}
	for name, value := range ms.counters {
		err = sendMetric(url, string(name), fmt.Sprintf("%d", value))
	}

	// return last error
	return err
}

func main() {
	pollDone := make(chan bool)
	reportDone := make(chan bool)
	pollQuit := make(chan bool)
	reportQuit := make(chan bool)
	storage := &MemStorage{
		gauges:   make(map[_const.GaugeMetricName]float64, 28),
		counters: make(map[_const.CounterMetricName]int64, 1),
	}

	pollInterval := 2
	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	go func() {
		for {
			select {
			case <-pollTicker.C:
				storage.writeMetrics()
			case <-pollQuit:
				// exit from reporting if no polling
				reportQuit <- true
				pollTicker.Stop()
				pollDone <- true
				return
			}
		}
	}()

	reportInterval := 10
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	go func() {
		for {
			select {
			case <-reportTicker.C:
				if err := storage.sendMetrics("http://localhost:8080"); err != nil {
					fmt.Println(err.Error())
				}
			case <-reportQuit:
				reportTicker.Stop()
				reportDone <- true
				return
			}
		}
	}()

	<-reportDone
	<-pollDone
}
