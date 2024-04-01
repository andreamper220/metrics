package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/andreamper220/metrics.git/internal/constants"
)

type MemStorage struct {
	counters map[constants.CounterMetricName]int64
	gauges   map[constants.GaugeMetricName]float64
}

func (ms *MemStorage) writeMetrics() {
	var mstats runtime.MemStats
	ms.gauges = map[constants.GaugeMetricName]float64{
		constants.Alloc:         float64(mstats.Alloc),
		constants.BuckHashSys:   float64(mstats.BuckHashSys),
		constants.Frees:         float64(mstats.Frees),
		constants.GcCPUFraction: mstats.GCCPUFraction,
		constants.GcSys:         float64(mstats.GCSys),
		constants.HeapAlloc:     float64(mstats.HeapAlloc),
		constants.HeapIdle:      float64(mstats.HeapIdle),
		constants.HeapInuse:     float64(mstats.HeapInuse),
		constants.HeapObjects:   float64(mstats.HeapObjects),
		constants.HeapReleased:  float64(mstats.HeapReleased),
		constants.HeapSys:       float64(mstats.HeapSys),
		constants.LastGc:        float64(mstats.LastGC),
		constants.Lookups:       float64(mstats.Lookups),
		constants.MemCacheInuse: float64(mstats.MCacheInuse),
		constants.MemCacheSys:   float64(mstats.MCacheSys),
		constants.MemSpanInuse:  float64(mstats.MSpanInuse),
		constants.MemSpanSys:    float64(mstats.MSpanSys),
		constants.MemAllocs:     float64(mstats.Mallocs),
		constants.NextGc:        float64(mstats.NextGC),
		constants.NumForcedGc:   float64(mstats.NumForcedGC),
		constants.NumGc:         float64(mstats.NumGC),
		constants.OtherSys:      float64(mstats.OtherSys),
		constants.PauseTotalNs:  float64(mstats.PauseTotalNs),
		constants.RandomValue:   rand.Float64(),
		constants.StackInuse:    float64(mstats.StackInuse),
		constants.StackSys:      float64(mstats.StackSys),
		constants.Sys:           float64(mstats.Sys),
		constants.TotalAlloc:    float64(mstats.TotalAlloc),
	}
	ms.counters[constants.PollCount] = 1
}

func sendMetric(url, name, value string) error {
	requestURL := fmt.Sprintf("%s/update/%s/%s/%s", url, constants.GaugeMetricType, name, value)
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	return res.Body.Close()
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
		gauges:   make(map[constants.GaugeMetricName]float64, 28),
		counters: make(map[constants.CounterMetricName]int64, 1),
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
