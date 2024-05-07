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
	SetToSaveMetricsAsync(bool) error
}

type metrics struct {
	counters map[shared.CounterMetricName]int64
	gauges   map[shared.GaugeMetricName]float64
}

type AbstractStorage struct {
	metrics            metrics
	toSaveMetricsAsync bool
}

func (as *AbstractStorage) GetCounters() map[shared.CounterMetricName]int64 {
	return as.metrics.counters
}
func (as *AbstractStorage) SetCounters(counters map[shared.CounterMetricName]int64) error {
	for name, value := range counters {
		as.metrics.counters[name] = value
	}
	return nil
}
func (as *AbstractStorage) GetGauges() map[shared.GaugeMetricName]float64 {
	return as.metrics.gauges
}
func (as *AbstractStorage) SetGauges(gauges map[shared.GaugeMetricName]float64) error {
	for name, value := range gauges {
		as.metrics.gauges[name] = value
	}
	return nil
}
func (as *AbstractStorage) GetToSaveMetricsAsync() bool {
	return as.toSaveMetricsAsync
}
func (as *AbstractStorage) SetToSaveMetricsAsync(toSaveMetricsAsync bool) error {
	as.toSaveMetricsAsync = toSaveMetricsAsync
	return nil
}
func NewAbstractStorage() *AbstractStorage {
	return &AbstractStorage{
		metrics: metrics{
			counters: make(map[shared.CounterMetricName]int64),
			gauges:   make(map[shared.GaugeMetricName]float64),
		},
	}
}
