package storages

import "github.com/andreamper220/metrics.git/internal/shared"

var Storage StorageInterface

// StorageInterface определяет набор методов для хранилища метрик.
type StorageInterface interface {
	GetCounters() ([]CounterMetric, error) // получение метрик-счётчиков
	AddCounter(CounterMetric) error        // добавление метрики-счётчика
	GetGauges() ([]GaugeMetric, error)     // получение метрик-значений
	AddGauge(GaugeMetric) error            // добавление метрики-значения
	GetMetrics() (Metrics, error)          // получение всех метрик
}

// CounterMetric определяет структуру метрики-счётчика для внутреннего использования.
type CounterMetric struct {
	Name  shared.CounterMetricName
	Value int64
}

// GaugeMetric определяет структуру метрики-значения для внутреннего использования.
type GaugeMetric struct {
	Name  shared.GaugeMetricName
	Value float64
}

// Metrics определяет структуру для хранения всех метрик.
type Metrics struct {
	counters []CounterMetric
	gauges   []GaugeMetric
}
