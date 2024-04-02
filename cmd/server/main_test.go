package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreamper220/metrics.git/internal/constants"
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
				constants.GaugeMetricType,
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
				constants.CounterMetricType,
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
				constants.GaugeMetricType,
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
				constants.GaugeMetricType,
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
				constants.GaugeMetricType,
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
				constants.CounterMetricType,
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
				r := chi.NewRouter()
				r.Post(`/update/{type}/{name}/{value}`, metric.UpdateMetric)
				server := httptest.NewServer(r)
				defer server.Close()

				req, _ := http.NewRequest(tt.got.method, fmt.Sprintf("%s/update/%s/%s/%s",
					server.URL, tt.got.metricType, tt.got.metricName, tt.got.metricValue), nil)
				req.Header.Set("Content-Type", "text/plain")

				client := &http.Client{}
				res, _ := client.Do(req)
				assert.Equal(t, tt.want.code, res.StatusCode)
				assert.NoError(t, res.Body.Close())
			})
	}
}

func TestShowMetric(t *testing.T) {
	type request struct {
		method     string
		metricType string
		metricName string
	}

	type response struct {
		code int
	}

	tests := []struct {
		request  request
		response response
		value    string
	}{
		{
			request: request{
				method:     http.MethodPost,
				metricType: constants.CounterMetricType,
				metricName: "invalid_method",
			},
			response: response{
				code: http.StatusMethodNotAllowed,
			},
			value: "",
		},
		{
			request: request{
				method:     http.MethodGet,
				metricType: constants.CounterMetricType,
				metricName: "valid_counter",
			},
			response: response{
				code: http.StatusOK,
			},
			value: "1",
		},
		{
			request: request{
				method:     http.MethodGet,
				metricType: constants.GaugeMetricType,
				metricName: "valid_gauge",
			},
			response: response{
				code: http.StatusOK,
			},
			value: "1.5",
		},
		{
			request: request{
				method:     http.MethodGet,
				metricType: "type",
				metricName: "invalid_type",
			},
			response: response{
				code: http.StatusBadRequest,
			},
			value: "",
		},
		{
			request: request{
				method:     http.MethodGet,
				metricType: constants.CounterMetricType,
				metricName: "not_found_counter",
			},
			response: response{
				code: http.StatusNotFound,
			},
			value: "",
		},
		{
			request: request{
				method:     http.MethodGet,
				metricType: constants.GaugeMetricType,
				metricName: "not_found_gauge",
			},
			response: response{
				code: http.StatusNotFound,
			},
			value: "",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s /value/%s/%s", tt.request.method, tt.request.metricType, tt.request.metricName),
			func(t *testing.T) {
				r := chi.NewRouter()
				r.Post(`/update/{type}/{name}/{value}`, metric.UpdateMetric)
				r.Get(`/value/{type}/{name}`, metric.ShowMetric)
				server := httptest.NewServer(r)
				defer server.Close()

				client := &http.Client{}
				if tt.value != "" {
					res, err := client.Post(fmt.Sprintf("%s/update/%s/%s/%s",
						server.URL, tt.request.metricType, tt.request.metricName, tt.value), "text/plain", nil)
					require.NoError(t, err)
					require.NoError(t, res.Body.Close())
				}

				req, _ := http.NewRequest(tt.request.method, fmt.Sprintf("%s/value/%s/%s",
					server.URL, tt.request.metricType, tt.request.metricName), nil)
				res, _ := client.Do(req)
				assert.Equal(t, tt.response.code, res.StatusCode)

				if tt.value != "" {
					resp, err := io.ReadAll(res.Body)
					require.NoError(t, err)
					assert.Equal(t, tt.value, string(resp))
				}
				require.NoError(t, res.Body.Close())
			})
	}
}

func TestShowMetrics(t *testing.T) {
	r := chi.NewRouter()
	r.Get(`/`, metric.ShowMetrics)
	server := httptest.NewServer(r)
	defer server.Close()

	res, err := http.Get(fmt.Sprintf(`%s/`, server.URL))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	resp, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.NotEmpty(t, string(resp))
	require.NoError(t, res.Body.Close())
}
