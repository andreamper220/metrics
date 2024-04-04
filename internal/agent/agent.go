package agent

import (
	"fmt"
	"time"

	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

func Run() error {
	blockDone := make(chan bool)

	storage := &storages.MemStorage{
		Gauges:   make(map[shared.GaugeMetricName]float64, 28),
		Counters: make(map[shared.CounterMetricName]int64, 1),
	}

	pollTicker := time.NewTicker(time.Duration(Config.PollInterval) * time.Second)
	go func() {
		for {
			select {
			case <-pollTicker.C:
				storage.WriteMetrics()
			}
		}
	}()

	reportTicker := time.NewTicker(time.Duration(Config.ReportInterval) * time.Second)
	go func() {
		for {
			select {
			case <-reportTicker.C:
				if err := storage.SendMetrics("http://" + Config.ServerAddress.String()); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}()

	<-blockDone

	return nil
}
