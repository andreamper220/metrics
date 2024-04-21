package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

func SendMetric(url string, metric shared.Metric, client *http.Client) error {
	requestURL := fmt.Sprintf("%s/update/", url)
	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	res, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return res.Body.Close()
}

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
			case <-blockDone:
				pollTicker.Stop()
				return
			}
		}
	}()

	reportTicker := time.NewTicker(time.Duration(Config.ReportInterval) * time.Second)
	go func() {
		for {
			select {
			case <-reportTicker.C:
				url := "http://" + Config.ServerAddress.String()
				client := &http.Client{
					Timeout: 30 * time.Second,
				}

				for name, value := range storage.Gauges {
					if err := SendMetric(url, shared.Metric{
						ID:    string(name),
						MType: shared.GaugeMetricType,
						Value: &value,
					}, client); err != nil {
						fmt.Println(err.Error())
					}
				}
				for name, value := range storage.Counters {
					if err := SendMetric(url, shared.Metric{
						ID:    string(name),
						MType: shared.CounterMetricType,
						Delta: &value,
					}, client); err != nil {
						fmt.Println(err.Error())
					}
				}
			case <-blockDone:
				reportTicker.Stop()
				return
			}
		}
	}()

	<-blockDone

	return nil
}
