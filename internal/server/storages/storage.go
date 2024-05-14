package storages

import "github.com/andreamper220/metrics.git/internal/shared"

var Storage StorageInterface

type StorageInterface interface {
	WriteMetrics() error
	ReadMetrics() error
	GetCounters() map[shared.CounterMetricName]int64
	SetCounters(map[shared.CounterMetricName]int64) error
	GetGauges() map[shared.GaugeMetricName]float64
	SetGauges(map[shared.GaugeMetricName]float64) error
	GetToSaveMetricsAsync() bool
}

type metrics struct {
	counters map[shared.CounterMetricName]int64
	gauges   map[shared.GaugeMetricName]float64
}
