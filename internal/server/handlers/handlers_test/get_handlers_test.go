package handlers_test

import (
	"fmt"
	"github.com/andreamper220/metrics.git/internal/shared"
	"io"
	"net/http"
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
					res, err := client.Post(fmt.Sprintf("%s/update/%s/%s/%s",
						s.Server.URL, tt.request.metricType, tt.request.metricName, tt.value), "text/plain", nil)
					s.Require().NoError(err)
					s.Require().NoError(res.Body.Close())
				}

				req, err := http.NewRequest(tt.request.method, fmt.Sprintf("%s/value/%s/%s",
					s.Server.URL, tt.request.metricType, tt.request.metricName), nil)
				s.Require().NoError(err)
				res, err := client.Do(req)
				s.Require().NoError(err)
				s.Equal(tt.response.code, res.StatusCode)

				if tt.value != "" {
					resp, err := io.ReadAll(res.Body)
					s.Require().NoError(err)
					s.Equal(tt.value, string(resp))
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
