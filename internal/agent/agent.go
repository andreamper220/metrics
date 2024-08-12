package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

type requestStruct struct {
	url        string
	bodyStruct interface{}
	client     *http.Client
}

func Run(requestCh chan requestStruct, errCh chan error) error {
	if err := logger.Initialize(); err != nil {
		return err
	}

	go updateMetrics()
	go updatePsUtilsMetrics()

	serverless := true
	if requestCh == nil && errCh == nil {
		requestCh = make(chan requestStruct)
		errCh = make(chan error)
		serverless = false
	}
	for s := 1; s <= Config.RateLimit; s++ {
		go Sender(requestCh, errCh)
	}
	go sendMetrics(requestCh)

	if !serverless {
		for err := range errCh {
			logger.Log.Error(err.Error())
		}
	}

	return nil
}

func Sender(requestCh <-chan requestStruct, errCh chan<- error) {
	defer close(errCh)

	for request := range requestCh {
		body, err := json.Marshal(request.bodyStruct)
		if err != nil {
			errCh <- err
			continue
		}

		// gzip compression
		var b bytes.Buffer
		zw := gzip.NewWriter(&b)
		if _, err = zw.Write(body); err != nil {
			errCh <- err
			continue
		}
		if err = zw.Close(); err != nil {
			errCh <- err
			continue
		}

		// hmac sha256
		var hash []byte
		if Config.Sha256Key != "" {
			h := hmac.New(sha256.New, []byte(Config.Sha256Key))
			if _, err = h.Write(body); err != nil {
				errCh <- err
				continue
			}
			hash = h.Sum(nil)
		}

		err = retry.Do(
			func() error {
				req, _ := http.NewRequest(http.MethodPost, request.url, &b)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Encoding", "gzip")
				if hash != nil {
					req.Header.Set("Hash-Sha256", hex.EncodeToString(hash))
				}
				res, err := request.client.Do(req)
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

		if err != nil {
			errCh <- err
			continue
		}
	}
}

func sendMetrics(requestCh chan<- requestStruct) {
	defer close(requestCh)

	reportTicker := time.NewTicker(time.Duration(Config.ReportInterval) * time.Second)
	for range reportTicker.C {
		currentMetrics := readMetrics()

		url := "http://" + Config.ServerAddress.String() + "/updates/"
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		metrics := make([]shared.Metric, len(currentMetrics.Gauges)+len(currentMetrics.Counters))
		metricsIndex := 0
		for name, value := range currentMetrics.Gauges {
			metrics[metricsIndex] = shared.Metric{
				ID:    string(name),
				MType: shared.GaugeMetricType,
				Value: &value,
			}
			metricsIndex++
		}
		for name, value := range currentMetrics.Counters {
			metrics[metricsIndex] = shared.Metric{
				ID:    string(name),
				MType: shared.CounterMetricType,
				Delta: &value,
			}
			metricsIndex++
		}

		requestCh <- requestStruct{
			url:        url,
			bodyStruct: metrics,
			client:     client,
		}
	}
}
