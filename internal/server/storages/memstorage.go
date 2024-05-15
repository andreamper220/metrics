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
	for _, counter := range ms.metrics.counters {
		if counter.Name == metric.Name {
			counter.Value = metric.Value
			isExisted = true
			break
		}
	}
	if !isExisted {
		ms.metrics.counters = append(ms.metrics.counters, CounterMetric{
			Name:  metric.Name,
			Value: metric.Value,
		})
	}

	return nil
}
func (ms *MemStorage) AddCounters(metrics []CounterMetric) error {
	var err error
	for _, metric := range metrics {
		err = ms.AddCounter(metric)
	}
	return err
}
func (ms *MemStorage) GetGauges() ([]GaugeMetric, error) {
	return ms.metrics.gauges, nil
}
func (ms *MemStorage) AddGauge(metric GaugeMetric) error {
	isExisted := false
	for _, gauge := range ms.metrics.gauges {
		if gauge.Name == metric.Name {
			gauge.Value = metric.Value
			isExisted = true
			break
		}
	}
	if !isExisted {
		ms.metrics.gauges = append(ms.metrics.gauges, GaugeMetric{
			Name:  metric.Name,
			Value: metric.Value,
		})
	}

	return nil
}
func (ms *MemStorage) AddGauges(metrics []GaugeMetric) error {
	var err error
	for _, metric := range metrics {
		err = ms.AddGauge(metric)
	}
	return err
}
func (ms *MemStorage) GetMetrics() (Metrics, error) {
	return ms.metrics, nil
}
