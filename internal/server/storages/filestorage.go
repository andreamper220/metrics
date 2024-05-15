package storages

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/shared"
)

type FileStorage struct {
	metrics            Metrics
	fileStoragePath    string
	toSaveMetricsAsync bool
}

func NewFileStorage(
	fileStoragePath string, storeInterval int, toRestore bool, blockDone chan bool,
) (*FileStorage, error) {
	fs := &FileStorage{
		metrics:            Metrics{},
		fileStoragePath:    fileStoragePath,
		toSaveMetricsAsync: storeInterval == 0,
	}

	// to restore metrics from file
	if toRestore {
		if err := fs.ReadMetrics(); err != nil {
			return nil, err
		}
	}
	// to store metrics to file
	if storeInterval > 0 {
		storeTicker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		go func() {
			for {
				select {
				case <-storeTicker.C:
					if err := fs.WriteMetrics(); err != nil {
						logger.Log.Error(err.Error())
					}
				case <-blockDone:
					storeTicker.Stop()
					return
				}
			}
		}()
	}
	return fs, nil
}
func (fs *FileStorage) GetCounters() ([]CounterMetric, error) {
	return fs.metrics.counters, nil
}
func (fs *FileStorage) AddCounter(metric CounterMetric) error {
	isExisted := false
	for _, counter := range fs.metrics.counters {
		if counter.Name == metric.Name {
			counter.Value = metric.Value
			isExisted = true
			break
		}
	}
	if !isExisted {
		fs.metrics.counters = append(fs.metrics.counters, CounterMetric{
			Name:  metric.Name,
			Value: metric.Value,
		})
	}

	if fs.toSaveMetricsAsync {
		return fs.WriteMetrics()
	}

	return nil
}
func (fs *FileStorage) AddCounters(metrics []CounterMetric) error {
	var err error
	for _, metric := range metrics {
		err = fs.AddCounter(metric)
	}

	if fs.toSaveMetricsAsync {
		return fs.WriteMetrics()
	}

	return err
}
func (fs *FileStorage) GetGauges() ([]GaugeMetric, error) {
	return fs.metrics.gauges, nil
}
func (fs *FileStorage) AddGauge(metric GaugeMetric) error {
	isExisted := false
	for _, gauge := range fs.metrics.gauges {
		if gauge.Name == metric.Name {
			gauge.Value = metric.Value
			isExisted = true
			break
		}
	}
	if !isExisted {
		fs.metrics.gauges = append(fs.metrics.gauges, GaugeMetric{
			Name:  metric.Name,
			Value: metric.Value,
		})
	}

	if fs.toSaveMetricsAsync {
		return fs.WriteMetrics()
	}

	return nil
}
func (fs *FileStorage) AddGauges(metrics []GaugeMetric) error {
	var err error
	for _, metric := range metrics {
		err = fs.AddGauge(metric)
	}

	if fs.toSaveMetricsAsync {
		return fs.WriteMetrics()
	}

	return err
}
func (fs *FileStorage) GetMetrics() (Metrics, error) {
	return fs.metrics, nil
}
func (fs *FileStorage) WriteMetrics() error {
	file, err := os.OpenFile(fs.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	data := make([]byte, 0)
	for _, counterMetric := range fs.metrics.counters {
		metric := &shared.Metric{
			ID:    string(counterMetric.Name),
			MType: shared.CounterMetricType,
			Delta: &counterMetric.Value,
		}
		metricData, _ := json.Marshal(&metric)
		metricData = append(metricData, '\n')
		data = append(data, metricData...)
	}

	for _, gaugeMetric := range fs.metrics.gauges {
		metric := &shared.Metric{
			ID:    string(gaugeMetric.Name),
			MType: shared.GaugeMetricType,
			Value: &gaugeMetric.Value,
		}
		metricData, _ := json.Marshal(&metric)
		metricData = append(metricData, '\n')
		data = append(data, metricData...)
	}

	_, err = file.Write(data)

	return err
}
func (fs *FileStorage) ReadMetrics() error {
	file, err := os.OpenFile(fs.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	fr := bufio.NewReader(file)
	dec := json.NewDecoder(fr)
	for {
		var metric shared.Metric

		err := dec.Decode(&metric)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch metric.MType {
		case shared.CounterMetricType:
			fs.AddCounter(CounterMetric{
				Name:  shared.CounterMetricName(metric.ID),
				Value: *metric.Delta,
			})
		case shared.GaugeMetricType:
			fs.AddGauge(GaugeMetric{
				Name:  shared.GaugeMetricName(metric.ID),
				Value: *metric.Value,
			})
		default:
			logger.Log.Fatalf("incorrect metric: %s", metric.ID)
		}
	}

	return err
}
