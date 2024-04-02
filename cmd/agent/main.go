package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	config "github.com/andreamper220/metrics.git/internal/config/agent"
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

func sendMetric(url, mType, name, value string, client *http.Client) error {
	requestURL := fmt.Sprintf("%s/update/%s/%s/%s", url, mType, name, value)

	res, err := client.Post(requestURL, "text/plain", http.NoBody)
	if err != nil {
		return err
	}

	return res.Body.Close()
}

func (ms *MemStorage) sendMetrics(url string) error {
	var err error = nil
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for name, value := range ms.gauges {
		err = sendMetric(url, constants.GaugeMetricType, string(name), fmt.Sprintf("%f", value), client)
	}
	for name, value := range ms.counters {
		err = sendMetric(url, constants.CounterMetricType, string(name), fmt.Sprintf("%d", value), client)
	}

	// return last error
	return err
}

func main() {
	config.ParseFlags()
	fmt.Println(config.Config.ServerAddress.String())
	fmt.Println(config.Config.ReportInterval)
	fmt.Println(config.Config.PollInterval)

	pollDone := make(chan bool)
	reportDone := make(chan bool)
	pollQuit := make(chan bool)
	reportQuit := make(chan bool)
	storage := &MemStorage{
		gauges:   make(map[constants.GaugeMetricName]float64, 28),
		counters: make(map[constants.CounterMetricName]int64, 1),
	}

	pollTicker := time.NewTicker(time.Duration(config.Config.PollInterval) * time.Second)
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

	reportTicker := time.NewTicker(time.Duration(config.Config.ReportInterval) * time.Second)
	go func() {
		for {
			select {
			case <-reportTicker.C:
				if err := storage.sendMetrics("http://" + config.Config.ServerAddress.String()); err != nil {
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
