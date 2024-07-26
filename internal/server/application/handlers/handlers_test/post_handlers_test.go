package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"golang.org/x/exp/constraints"

	"github.com/andreamper220/metrics.git/internal/server/domain/metrics"
	"github.com/andreamper220/metrics.git/internal/server/infrastructure/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func Ptr[T Number](v T) *T {
	return &v
}

func (s *HandlerTestSuite) TestUpdateMetric() {
	type request struct {
		method      string
		metricType  string
		metricName  string
		metricValue string
	}

	type response struct {
		code int
	}

	tests := []struct {
		got  request
		want response
	}{
		{
			request{
				http.MethodPost,
				shared.GaugeMetricType,
				"test_metric_gauge",
				"2.50",
			},
			response{
				http.StatusOK,
			},
		},
		{
			request{
				http.MethodPost,
				shared.CounterMetricType,
				"test_metric_counter",
				"2",
			},
			response{
				http.StatusOK,
			},
		},
		{
			request{
				http.MethodGet,
				shared.GaugeMetricType,
				"test_metric_gauge",
				"2.50",
			},
			response{
				http.StatusMethodNotAllowed,
			},
		},
		{
			request{
				http.MethodPost,
				shared.GaugeMetricType,
				"",
				"2.50",
			},
			response{
				http.StatusNotFound,
			},
		},
		{
			request{
				http.MethodPost,
				"type",
				"test_metric_gauge",
				"2.50",
			},
			response{
				http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /update/%s/%s/%s", tt.got.method, tt.got.metricType, tt.got.metricName, tt.got.metricValue),
			func() {
				metric := shared.Metric{
					ID:    tt.got.metricName,
					MType: tt.got.metricType,
				}
				switch tt.got.metricType {
				case shared.CounterMetricType:
					value, err := strconv.ParseInt(tt.got.metricValue, 10, 64)
					s.Require().NoError(err)
					metric.Delta = &value
				case shared.GaugeMetricType:
					value, err := strconv.ParseFloat(tt.got.metricValue, 64)
					s.Require().NoError(err)
					metric.Value = &value
				}
				body, err := json.Marshal(metric)
				s.Require().NoError(err)

				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/update/", s.Server.URL),
					bytes.NewBuffer(body),
				)
				s.Require().NoError(err)
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				res, err := client.Do(req)
				s.Require().NoError(err)
				s.Equal(tt.want.code, res.StatusCode)
				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}

func (s *HandlerTestSuite) TestUpdateMetrics() {
	type request struct {
		method             string
		metricNameGauge    string
		metricValueGauge   float64
		metricNameCounter  string
		metricValueCounter int64
	}

	type response struct {
		code int
	}

	tests := []struct {
		got  request
		want response
	}{
		{
			request{
				http.MethodPost,
				"test_metric_gauge",
				2.5,
				"test_metric_counter",
				17,
			},
			response{
				http.StatusOK,
			},
		},
		{
			request{
				http.MethodGet,
				"test_metric_gauge",
				2.5,
				"test_metric_counter",
				17,
			},
			response{
				http.StatusMethodNotAllowed,
			},
		},
	}

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /updates/", tt.got.method),
			func() {
				metrics := []shared.Metric{
					{
						ID:    tt.got.metricNameGauge,
						MType: shared.GaugeMetricType,
						Value: &tt.got.metricValueGauge,
					},
					{
						ID:    tt.got.metricNameCounter,
						MType: shared.CounterMetricType,
						Delta: &tt.got.metricValueCounter,
					},
				}
				body, err := json.Marshal(metrics)
				s.Require().NoError(err)

				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/updates/", s.Server.URL),
					bytes.NewBuffer(body),
				)
				s.Require().NoError(err)
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				res, err := client.Do(req)
				s.Require().NoError(err)
				s.Equal(tt.want.code, res.StatusCode)
				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}

func (s *HandlerTestSuite) TestUpdateMetricOld() {
	type request struct {
		method      string
		metricType  string
		metricName  string
		metricValue string
	}

	type response struct {
		code int
	}

	tests := []struct {
		got  request
		want response
	}{
		{
			request{
				http.MethodPost,
				shared.GaugeMetricType,
				"test_metric_gauge",
				"2.50",
			},
			response{
				http.StatusOK,
			},
		},
		{
			request{
				http.MethodPost,
				shared.CounterMetricType,
				"test_metric_counter",
				"2",
			},
			response{
				http.StatusOK,
			},
		},
		{
			request{
				http.MethodGet,
				shared.GaugeMetricType,
				"test_metric_gauge",
				"2.50",
			},
			response{
				http.StatusMethodNotAllowed,
			},
		},
		{
			request{
				http.MethodPost,
				shared.GaugeMetricType,
				"",
				"2.50",
			},
			response{
				http.StatusNotFound,
			},
		},
		{
			request{
				http.MethodPost,
				"type",
				"test_metric_gauge",
				"2.50",
			},
			response{
				http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /update/%s/%s/%s", tt.got.method, tt.got.metricType, tt.got.metricName, tt.got.metricValue),
			func() {
				req, err := http.NewRequest(
					tt.got.method,
					fmt.Sprintf("%s/update/%s/%s/%s", s.Server.URL, tt.got.metricType, tt.got.metricName, tt.got.metricValue),
					nil,
				)
				s.Require().NoError(err)

				client := &http.Client{}
				res, err := client.Do(req)
				s.Require().NoError(err)
				s.Equal(tt.want.code, res.StatusCode)
				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}

func BenchmarkProcessMetric(b *testing.B) {
	const metricsN = 100
	metricsSet := make([]shared.Metric, 0, metricsN)

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < metricsN; i++ {
		str := make([]rune, 10)
		for i := range str {
			str[i] = letters[rand.Intn(len(letters))]
		}
		metricsSet = append(metricsSet, shared.Metric{
			ID:    string(str),
			MType: shared.CounterMetricType,
			Delta: Ptr(int64(2)),
		})
	}
	storages.Storage = storages.NewMemStorage()

	b.ResetTimer()

	b.Run("metrics", func(b *testing.B) {
		for i := 0; i < metricsN; i++ {
			err := metrics.ProcessMetric(&metricsSet[i])
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
