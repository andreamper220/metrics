package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

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

	if reqMetric.ID == "" {
		http.Error(w, "Not found metric ID.", http.StatusNotFound)
		return
	}

	switch reqMetric.MType {
	case shared.CounterMetricType:
		var value = storages.Storage.GetCounters()[shared.CounterMetricName(reqMetric.ID)] + *reqMetric.Delta

		*reqMetric.Delta = value
		if err := storages.Storage.SetCounters(map[shared.CounterMetricName]int64{
			shared.CounterMetricName(reqMetric.ID): value,
		}); err != nil {
			logger.Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case shared.GaugeMetricType:
		if err := storages.Storage.SetGauges(map[shared.GaugeMetricName]float64{
			shared.GaugeMetricName(reqMetric.ID): *reqMetric.Value,
		}); err != nil {
			logger.Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Incorrect metric TYPE.", http.StatusBadRequest)
		return
	}

	if storages.Storage.GetToSaveMetricsAsync() {
		if err := storages.Storage.WriteMetrics(); err != nil {
			logger.Log.Error(err.Error())
		}
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

func UpdateMetricOld(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "Not found metric NAME.", http.StatusNotFound)
		return
	}

	switch chi.URLParam(r, "type") {
	case shared.CounterMetricType:
		value, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
		if err != nil {
			http.Error(w, "Incorrect metric VALUE: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := storages.Storage.SetCounters(map[shared.CounterMetricName]int64{
			shared.CounterMetricName(name): storages.Storage.GetCounters()[shared.CounterMetricName(name)] + value,
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

		if err := storages.Storage.SetGauges(map[shared.GaugeMetricName]float64{
			shared.GaugeMetricName(name): value,
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
