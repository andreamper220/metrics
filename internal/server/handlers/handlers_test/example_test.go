package handlers_test

import (
	"github.com/andreamper220/metrics.git/internal/server/handlers"
	"github.com/andreamper220/metrics.git/internal/shared"
)

func ExampleProcessMetric() {
	reqMetric := shared.Metric{
		ID:    "test_id",
		MType: shared.CounterMetricType,
		Delta: Ptr(int64(1)),
	}

	if err := handlers.ProcessMetric(&reqMetric); err != nil {
		return
	}
}
