package storages

import "github.com/andreamper220/metrics.git/internal/shared"

var Storage StorageInterface

type StorageInterface interface {
	GetCounters() ([]CounterMetric, error)
	AddCounter(CounterMetric) error
	AddCounters([]CounterMetric) error
	GetGauges() ([]GaugeMetric, error)
	AddGauge(GaugeMetric) error
	AddGauges([]GaugeMetric) error
	GetMetrics() (Metrics, error)
}

type CounterMetric struct {
	Name  shared.CounterMetricName
	Value int64
}

type GaugeMetric struct {
	Name  shared.GaugeMetricName
	Value float64
}

type Metrics struct {
	counters []CounterMetric
	gauges   []GaugeMetric
}
