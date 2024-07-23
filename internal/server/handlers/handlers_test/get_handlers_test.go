package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/andreamper220/metrics.git/internal/shared"
)

func (s *HandlerTestSuite) TestShowMetric() {
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
				method:     http.MethodGet,
				metricType: shared.CounterMetricType,
				metricName: "invalid_method",
			},
			response: response{
				code: http.StatusMethodNotAllowed,
			},
			value: "",
		},
		{
			request: request{
				method:     http.MethodPost,
				metricType: shared.CounterMetricType,
				metricName: "valid_counter",
			},
			response: response{
				code: http.StatusOK,
			},
			value: "1",
		},
		{
			request: request{
				method:     http.MethodPost,
				metricType: shared.GaugeMetricType,
				metricName: "valid_gauge",
			},
			response: response{
				code: http.StatusOK,
			},
			value: "1.5",
		},
		{
			request: request{
				method:     http.MethodPost,
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
				method:     http.MethodPost,
				metricType: shared.CounterMetricType,
				metricName: "not_found_counter",
			},
			response: response{
				code: http.StatusNotFound,
			},
			value: "",
		},
		{
			request: request{
				method:     http.MethodPost,
				metricType: shared.GaugeMetricType,
				metricName: "not_found_gauge",
			},
			response: response{
				code: http.StatusNotFound,
			},
			value: "",
		},
	}

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /value/%s/%s", tt.request.method, tt.request.metricType, tt.request.metricName),
			func() {
				var resMetric shared.Metric
				postMetric := shared.Metric{
					ID:    tt.request.metricName,
					MType: tt.request.metricType,
				}
				getMetric := postMetric

				client := &http.Client{}
				if tt.value != "" {
					switch tt.request.metricType {
					case shared.CounterMetricType:
						value, err := strconv.ParseInt(tt.value, 10, 64)
						s.Require().NoError(err)
						postMetric.Delta = &value
					case shared.GaugeMetricType:
						value, err := strconv.ParseFloat(tt.value, 64)
						s.Require().NoError(err)
						postMetric.Value = &value
					}
					body, err := json.Marshal(postMetric)
					s.Require().NoError(err)

					res, err := client.Post(
						fmt.Sprintf("%s/update/", s.Server.URL),
						"application/json",
						bytes.NewBuffer(body),
					)
					s.Require().NoError(err)
					s.Require().NoError(res.Body.Close())
				}

				body, err := json.Marshal(getMetric)
				s.Require().NoError(err)
				req, err := http.NewRequest(
					tt.request.method,
					fmt.Sprintf("%s/value/", s.Server.URL),
					bytes.NewBuffer(body),
				)
				s.Require().NoError(err)
				res, err := client.Do(req)
				s.Require().NoError(err)
				s.Equal(tt.response.code, res.StatusCode)

				if tt.value != "" {
					s.Require().NoError(json.NewDecoder(res.Body).Decode(&resMetric))
					var value string
					switch tt.request.metricType {
					case shared.CounterMetricType:
						value = fmt.Sprintf("%d", *resMetric.Delta)
					case shared.GaugeMetricType:
						value = fmt.Sprintf("%g", *resMetric.Value)
					}
					s.Equal(tt.value, value)
				}
				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}

func (s *HandlerTestSuite) TestShowMetrics() {
	res, err := http.Get(fmt.Sprintf(`%s/`, s.Server.URL))
	s.Require().NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
	resp, err := io.ReadAll(res.Body)
	s.Require().NoError(err)
	s.NotEmpty(string(resp))
	s.Require().NoError(res.Body.Close())

	defer s.Server.Close()
}

func (s *HandlerTestSuite) TestShowMetricOld() {
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
				metricType: shared.CounterMetricType,
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
				metricType: shared.CounterMetricType,
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
				metricType: shared.GaugeMetricType,
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
				metricType: shared.CounterMetricType,
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
				metricType: shared.GaugeMetricType,
				metricName: "not_found_gauge",
			},
			response: response{
				code: http.StatusNotFound,
			},
			value: "",
		},
	}

	for _, tt := range tests {
		s.Run(fmt.Sprintf("%s /value/%s/%s", tt.request.method, tt.request.metricType, tt.request.metricName),
			func() {
				client := &http.Client{}
				if tt.value != "" {
					req, err := http.NewRequest(
						http.MethodPost,
						fmt.Sprintf("%s/update/%s/%s/%s", s.Server.URL, tt.request.metricType, tt.request.metricName, tt.value),
						nil,
					)
					s.Require().NoError(err)

					res, err := client.Do(req)
					s.Require().NoError(err)
					s.Require().NoError(res.Body.Close())
				}

				req, err := http.NewRequest(
					tt.request.method,
					fmt.Sprintf("%s/value/%s/%s", s.Server.URL, tt.request.metricType, tt.request.metricName),
					nil,
				)
				s.Require().NoError(err)
				res, err := client.Do(req)
				s.Require().NoError(err)
				s.Equal(tt.response.code, res.StatusCode)

				if tt.value != "" {
					var buf bytes.Buffer
					_, err = buf.ReadFrom(res.Body)
					s.Require().NoError(err)
					s.Equal(tt.value, buf.String())
				}
				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}
