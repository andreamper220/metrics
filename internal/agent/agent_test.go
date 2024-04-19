package agent

import (
	"github.com/andreamper220/metrics.git/internal/logger"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreamper220/metrics.git/internal/server"
	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

func TestSendMetrics(t *testing.T) {
	tests := []struct {
		name    string
		storage storages.MemStorage
	}{
		{
			name: "send counter",
			storage: storages.MemStorage{
				Counters: map[shared.CounterMetricName]int64{
					shared.PollCount: 1,
				},
			},
		},
		{
			name: "send gauge",
			storage: storages.MemStorage{
				Gauges: map[shared.GaugeMetricName]float64{
					shared.Alloc: 2.5,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := logger.Initialize(); err != nil {
				t.Fatal(err.Error())
			}

			r := server.MakeRouter()
			srv := httptest.NewServer(r)
			defer srv.Close()

			assert.NoError(t, tt.storage.SendMetrics(srv.URL))
		})
	}
}
