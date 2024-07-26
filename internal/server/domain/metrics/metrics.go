package metrics

import (
	"errors"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/server/infrastructure/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

var (
	ErrMetricNotFound      = errors.New("not found metric ID")
	ErrIncorrectMetricType = errors.New("incorrect metric TYPE")
)

func ProcessMetric(metric *shared.Metric) error {
	if metric.ID == "" {
		return ErrMetricNotFound
	}

	switch metric.MType {
	case shared.CounterMetricType:
		counters, err := storages.Storage.GetCounters()
		if err != nil {
			return err
		}

		var value int64 = 0
		for _, counter := range counters {
			if counter.Name == shared.CounterMetricName(metric.ID) {
				value = counter.Value
				break
			}
		}
		value += *metric.Delta

		*metric.Delta = value
		if err := storages.Storage.AddCounter(storages.CounterMetric{
			Name:  shared.CounterMetricName(metric.ID),
			Value: value,
		}); err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	case shared.GaugeMetricType:
		if err := storages.Storage.AddGauge(storages.GaugeMetric{
			Name:  shared.GaugeMetricName(metric.ID),
			Value: *metric.Value,
		}); err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	default:
		return ErrIncorrectMetricType
	}

	return nil
}
