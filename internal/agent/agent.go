package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	ctxSignal, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	stopCh := make(chan struct{})
	go sendMetrics(ctxSignal, requestCh, stopCh)

	if !serverless {
		for {
			select {
			case err := <-errCh:
				logger.Log.Error(err.Error())
			case <-stopCh:
				return nil
			}
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

		// crypto
		if Config.CryptoKeyPath != "" {
			publicKeyPEM, err := os.ReadFile(Config.CryptoKeyPath)
			if err != nil {
				errCh <- err
				continue
			}
			publicKeyBlock, _ := pem.Decode(publicKeyPEM)
			publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
			if err != nil {
				errCh <- err
				continue
			}

			msgLen := len(body)
			step := publicKey.(*rsa.PublicKey).Size() / 2
			var encryptedBytes []byte

			for start := 0; start < msgLen; start += step {
				finish := start + step
				if finish > msgLen {
					finish = msgLen
				}

				cipherBody, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), body[start:finish])
				if err != nil {
					errCh <- err
					continue
				}

				encryptedBytes = append(encryptedBytes, cipherBody...)
			}
			b.Write(encryptedBytes)
		}

		err = retry.Do(
			func() error {
				req, _ := http.NewRequest(http.MethodPost, request.url, &b)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Content-Encoding", "gzip")
				if hash != nil {
					req.Header.Set("Hash-Sha256", hex.EncodeToString(hash))
				}
				req.Header.Set("X-Real-Ip", req.Host)
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

func sendMetrics(context context.Context, requestCh chan<- requestStruct, stopCh chan<- struct{}) {
	defer close(requestCh)

	reportTicker := time.NewTicker(time.Duration(Config.ReportInterval) * time.Second)
	for {
		select {
		case <-context.Done():
			reportTicker.Stop()
			stopCh <- struct{}{}
		case <-reportTicker.C:
			go func() {
				url := "http://" + Config.ServerAddress.String() + "/updates/"
				client := &http.Client{
					Timeout: 30 * time.Second,
				}

				requestCh <- requestStruct{
					url:        url,
					bodyStruct: buildMetrics(),
					client:     client,
				}
			}()
		}
	}
}

func buildMetrics() []shared.Metric {
	currentMetrics := readMetrics()

	metricsSlice := make([]shared.Metric, len(currentMetrics.Gauges)+len(currentMetrics.Counters))
	metricsIndex := 0
	for name, value := range currentMetrics.Gauges {
		metricsSlice[metricsIndex] = shared.Metric{
			ID:    string(name),
			MType: shared.GaugeMetricType,
			Value: &value,
		}
		metricsIndex++
	}
	for name, value := range currentMetrics.Counters {
		metricsSlice[metricsIndex] = shared.Metric{
			ID:    string(name),
			MType: shared.CounterMetricType,
			Delta: &value,
		}
		metricsIndex++
	}

	return metricsSlice
}
