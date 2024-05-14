package storages

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/shared"
)

type FileStorage struct {
	metrics            metrics
	toSaveMetricsAsync bool
	fileStoragePath    string
}

func NewFileStorage(fileStoragePath string, toSaveMetricsAsync bool) *FileStorage {
	return &FileStorage{
		metrics: metrics{
			counters: make(map[shared.CounterMetricName]int64),
			gauges:   make(map[shared.GaugeMetricName]float64),
		},
		toSaveMetricsAsync: toSaveMetricsAsync,
		fileStoragePath:    fileStoragePath,
	}
}
func (fs *FileStorage) GetCounters() map[shared.CounterMetricName]int64 {
	return fs.metrics.counters
}
func (fs *FileStorage) SetCounters(counters map[shared.CounterMetricName]int64) error {
	for name, value := range counters {
		fs.metrics.counters[name] = value
	}
	return nil
}
func (fs *FileStorage) GetGauges() map[shared.GaugeMetricName]float64 {
	return fs.metrics.gauges
}
func (fs *FileStorage) SetGauges(gauges map[shared.GaugeMetricName]float64) error {
	for name, value := range gauges {
		fs.metrics.gauges[name] = value
	}
	return nil
}
func (fs *FileStorage) GetToSaveMetricsAsync() bool {
	return fs.toSaveMetricsAsync
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
	for name, value := range fs.metrics.counters {
		metric := &shared.Metric{
			ID:    string(name),
			MType: shared.CounterMetricType,
			Delta: &value,
		}
		metricData, _ := json.Marshal(&metric)
		metricData = append(metricData, '\n')
		data = append(data, metricData...)
	}

	for name, value := range fs.metrics.gauges {
		metric := &shared.Metric{
			ID:    string(name),
			MType: shared.GaugeMetricType,
			Value: &value,
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
			fs.metrics.counters[shared.CounterMetricName(metric.ID)] = *metric.Delta
		case shared.GaugeMetricType:
			fs.metrics.gauges[shared.GaugeMetricName(metric.ID)] = *metric.Value
		default:
			logger.Log.Fatalf("incorrect metric: %s", metric.ID)
		}
	}

	return err
}
