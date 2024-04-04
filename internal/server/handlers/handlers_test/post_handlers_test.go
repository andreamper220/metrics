package handlers_test

import (
	"fmt"
	"net/http"

	"github.com/andreamper220/metrics.git/internal/shared"
)

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
		{
			request{
				http.MethodPost,
				shared.GaugeMetricType,
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
				shared.CounterMetricType,
				"test_metric_counter",
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
				req, err := http.NewRequest(tt.got.method, fmt.Sprintf("%s/update/%s/%s/%s",
					s.Server.URL, tt.got.metricType, tt.got.metricName, tt.got.metricValue), nil)
				s.Require().NoError(err)
				req.Header.Set("Content-Type", "text/plain")

				client := &http.Client{}
				res, err := client.Do(req)
				s.Require().NoError(err)
				s.Equal(tt.want.code, res.StatusCode)
				s.Require().NoError(res.Body.Close())
			})
	}

	defer s.Server.Close()
}
