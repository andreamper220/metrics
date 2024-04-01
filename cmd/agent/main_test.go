package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	constants "github.com/andreamper220/metrics.git/internal/const"
	metric "github.com/andreamper220/metrics.git/internal/server"
)

func TestSendMetrics(t *testing.T) {
	tests := []struct {
		name    string
		storage MemStorage
	}{
		{
			name: "send counter",
			storage: MemStorage{
				counters: map[constants.CounterMetricName]int64{
					constants.PollCount: 1,
				},
			},
		},
		{
			name: "send gauge",
			storage: MemStorage{
				gauges: map[constants.GaugeMetricName]float64{
					constants.Alloc: 2.5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(`/update/{type}/{name}/{value}`, metric.UpdateMetric)
			server := httptest.NewServer(mux)
			defer server.Close()

			assert.NoError(t, tt.storage.sendMetrics(server.URL))
		})
	}
}
