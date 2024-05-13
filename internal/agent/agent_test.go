package agent

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/server"
	"github.com/andreamper220/metrics.git/internal/shared"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func Ptr[T Number](v T) *T {
	return &v
}

func TestSendMetrics(t *testing.T) {
	tests := []struct {
		name   string
		metric shared.Metric
	}{
		{
			name: "send counter",
			metric: shared.Metric{
				ID:    string(shared.PollCount),
				MType: shared.CounterMetricType,
				Delta: Ptr(int64(1)),
			},
		},
		{
			name: "send gauge",
			metric: shared.Metric{
				ID:    string(shared.Alloc),
				MType: shared.GaugeMetricType,
				Value: Ptr(2.5),
			},
		},
	}

	for _, tt := range tests {
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := logger.Initialize(); err != nil {
				t.Fatal(err.Error())
			}
			if err := server.MakeStorage(); err != nil {
				t.Fatal(err.Error())
			}

			r := server.MakeRouter()
			srv := httptest.NewServer(r)
			defer srv.Close()

			require.NoError(t, Send(srv.URL+"/update/", tt.metric, client))
			body, err := json.Marshal(tt.metric)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, srv.URL+"/value/", bytes.NewBuffer(body))
			require.NoError(t, err)
			res, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			var resMetric shared.Metric
			require.NoError(t, json.NewDecoder(res.Body).Decode(&resMetric))
			switch tt.metric.MType {
			case shared.CounterMetricType:
				assert.Equal(t, tt.metric.Delta, resMetric.Delta)
			case shared.GaugeMetricType:
				assert.Equal(t, tt.metric.Value, resMetric.Value)
			}
			require.NoError(t, res.Body.Close())
		})
	}
}
