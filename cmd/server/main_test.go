package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	_const "github.com/andreamper220/metrics.git/internal/const"
	metric "github.com/andreamper220/metrics.git/internal/server"
)

func TestUpdateMetric(t *testing.T) {
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
				_const.GaugeMetricType,
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
				_const.CounterMetricType,
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
				_const.GaugeMetricType,
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
				_const.GaugeMetricType,
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
		{
			request{
				http.MethodPost,
				_const.GaugeMetricType,
				"test_metric_gauge",
				"metric",
			},
			response{
				http.StatusBadRequest,
			},
		},
		{
			request{
				http.MethodPost,
				_const.CounterMetricType,
				"test_metric_counter",
				"2.50",
			},
			response{
				http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s /update/%s/%s/%s", tt.got.method, tt.got.metricType, tt.got.metricName, tt.got.metricValue),
			func(t *testing.T) {
				mux := http.NewServeMux()
				mux.HandleFunc(`/update/{type}/{name}/{value}`, metric.UpdateMetric)
				server := httptest.NewServer(mux)
				defer server.Close()

				switch tt.got.method {
				case http.MethodPost:
					res, _ := http.Post(fmt.Sprintf("%s/update/%s/%s/%s",
						server.URL, tt.got.metricType, tt.got.metricName, tt.got.metricValue),
						"text/plain", nil,
					)
					assert.Equal(t, tt.want.code, res.StatusCode)
				case http.MethodGet:
					res, _ := http.Get(fmt.Sprintf("%s/update/%s/%s/%s",
						server.URL, tt.got.metricType, tt.got.metricName, tt.got.metricValue),
					)
					assert.Equal(t, tt.want.code, res.StatusCode)
				default:
					t.Fatal("No such method: " + tt.got.method)
				}
			})
	}
}
