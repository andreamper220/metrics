package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go"

	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

func SendMetric(url string, metric shared.Metric, client *http.Client) error {
	requestURL := fmt.Sprintf("%s/update/", url)
	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	err = retry.Do(
		func() error {
			res, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				var netErr net.Error
				if (errors.As(err, &netErr) && netErr.Timeout()) ||
					strings.Contains(err.Error(), "EOF") ||
					strings.Contains(err.Error(), "connection reset by peer") {
					return err // retry only network errors
				}
				return retry.Unrecoverable(err)
			}
			err = res.Body.Close()
			if err != nil {
				return retry.Unrecoverable(err)
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(time.Second),
		retry.DelayType(retry.BackOffDelay),
	)

	return err
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
