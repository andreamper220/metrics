package storages

type MemStorage struct {
	metrics Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Metrics{},
	}
}
func (ms *MemStorage) GetCounters() ([]CounterMetric, error) {
	return ms.metrics.counters, nil
}
func (ms *MemStorage) AddCounter(metric CounterMetric) error {
	isExisted := false
	for key, counter := range ms.metrics.counters {
		if counter.Name == metric.Name {
			ms.metrics.counters[key].Value = metric.Value
			isExisted = true
			break
		}
	}
	if !isExisted {
		ms.metrics.counters = append(ms.metrics.counters, metric)
	}

	return nil
}
func (ms *MemStorage) GetGauges() ([]GaugeMetric, error) {
	return ms.metrics.gauges, nil
}
func (ms *MemStorage) AddGauge(metric GaugeMetric) error {
	isExisted := false
	for key, gauge := range ms.metrics.gauges {
		if gauge.Name == metric.Name {
			ms.metrics.gauges[key].Value = metric.Value
			isExisted = true
			break
		}
	}
	if !isExisted {
		ms.metrics.gauges = append(ms.metrics.gauges, metric)
	}

	return nil
}
func (ms *MemStorage) GetMetrics() (Metrics, error) {
	return ms.metrics, nil
}
