package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

var (
	ErrMetricNotFound      = errors.New("not found metric ID")
	ErrIncorrectMetricType = errors.New("incorrect metric TYPE")
)

// UpdateMetric обновляет значение одной метрики, переданной в теле JSON.
func UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var reqMetric shared.Metric
	var buf bytes.Buffer

	// json unmarshal
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &reqMetric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := ProcessMetric(&reqMetric); err != nil {
		if errors.Is(err, ErrMetricNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if errors.Is(err, ErrIncorrectMetricType) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// json marshal
	resp, err := json.Marshal(reqMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateMetricOld обновляет значение искомой метрики.
//
// Deprecated: используйте UpdateMetric.
func UpdateMetricOld(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "Not found metric NAME.", http.StatusNotFound)
		return
	}

	switch chi.URLParam(r, "type") {
	case shared.CounterMetricType:
		metricValue, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
		if err != nil {
			http.Error(w, "Incorrect metric VALUE: "+err.Error(), http.StatusBadRequest)
			return
		}

		counters, err := storages.Storage.GetCounters()
		if err != nil {
			logger.Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var value int64 = 0
		for _, counter := range counters {
			if counter.Name == shared.CounterMetricName(name) {
				value = counter.Value
				break
			}
		}
		value += metricValue

		if err := storages.Storage.AddCounter(storages.CounterMetric{
			Name:  shared.CounterMetricName(name),
			Value: value,
		}); err != nil {
			logger.Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case shared.GaugeMetricType:
		value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
		if err != nil {
			http.Error(w, "Incorrect metric VALUE: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := storages.Storage.AddGauge(storages.GaugeMetric{
			Name:  shared.GaugeMetricName(name),
			Value: value,
		}); err != nil {
			logger.Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Incorrect metric TYPE.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateMetrics обновляет значение нескольких метрик, переданных в теле JSON.
func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	var reqMetrics, resMetrics shared.Metrics
	if err := json.NewDecoder(r.Body).Decode(&reqMetrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, reqMetric := range reqMetrics {
		if err := ProcessMetric(&reqMetric); err != nil {
			if errors.Is(err, ErrMetricNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else if errors.Is(err, ErrIncorrectMetricType) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		isExisted := false
		for key, resMetric := range resMetrics {
			if resMetric.ID == reqMetric.ID {
				resMetrics[key] = reqMetric
				isExisted = true
			}
		}
		if !isExisted {
			resMetrics = append(resMetrics, reqMetric)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resMetrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
