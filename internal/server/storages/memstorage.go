package storages

import "github.com/andreamper220/metrics.git/internal/shared"

type MemStorage struct {
	metrics            metrics
	toSaveMetricsAsync bool
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: metrics{
			counters: make(map[shared.CounterMetricName]int64),
			gauges:   make(map[shared.GaugeMetricName]float64),
		},
		toSaveMetricsAsync: true,
	}
}
func (ms *MemStorage) GetCounters() map[shared.CounterMetricName]int64 {
	return ms.metrics.counters
}
func (ms *MemStorage) SetCounters(counters map[shared.CounterMetricName]int64) error {
	for name, value := range counters {
		ms.metrics.counters[name] = value
	}
	return nil
}
func (ms *MemStorage) GetGauges() map[shared.GaugeMetricName]float64 {
	return ms.metrics.gauges
}
func (ms *MemStorage) SetGauges(gauges map[shared.GaugeMetricName]float64) error {
	for name, value := range gauges {
		ms.metrics.gauges[name] = value
	}
	return nil
}
func (ms *MemStorage) GetToSaveMetricsAsync() bool {
	return ms.toSaveMetricsAsync
}
func (ms *MemStorage) WriteMetrics() error {
	return nil
}
func (ms *MemStorage) ReadMetrics() error {
	return nil
}
