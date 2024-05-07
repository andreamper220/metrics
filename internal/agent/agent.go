package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/shared"
)

func SendMetric(url string, metric shared.Metric, client *http.Client) error {
	requestURL := url + "/update/"
	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	// gzip compression
	var b bytes.Buffer
	zw := gzip.NewWriter(&b)
	if _, err := zw.Write(body); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}

	err = retry.Do(
		func() error {
			req, _ := http.NewRequest(http.MethodPost, requestURL, &b)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Encoding", "gzip")
			res, err := client.Do(req)
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

	pollTicker := time.NewTicker(time.Duration(Config.PollInterval) * time.Second)
	go func() {
		for {
			select {
			case <-pollTicker.C:
				ReadMetrics()
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

				for name, value := range Metrics.Gauges {
					if err := SendMetric(url, shared.Metric{
						ID:    string(name),
						MType: shared.GaugeMetricType,
						Value: &value,
					}, client); err != nil {
						logger.Log.Error(err.Error())
					}
				}
				for name, value := range Metrics.Counters {
					if err := SendMetric(url, shared.Metric{
						ID:    string(name),
						MType: shared.CounterMetricType,
						Delta: &value,
					}, client); err != nil {
						logger.Log.Error(err.Error())
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
