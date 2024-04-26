package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

var storage = &storages.MemStorage{
	Counters: make(map[shared.CounterMetricName]int64),
	Gauges:   make(map[shared.GaugeMetricName]float64),
}

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
		storage.Counters[shared.CounterMetricName(reqMetric.ID)] = *reqMetric.Delta
	case shared.GaugeMetricType:
		storage.Gauges[shared.GaugeMetricName(reqMetric.ID)] = *reqMetric.Value
	default:
		http.Error(w, "Incorrect metric TYPE.", http.StatusBadRequest)
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

		storage.Counters[shared.CounterMetricName(name)] = value
	case shared.GaugeMetricType:
		value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
		if err != nil {
			http.Error(w, "Incorrect metric VALUE: "+err.Error(), http.StatusBadRequest)
			return
		}

		storage.Gauges[shared.GaugeMetricName(name)] = value
	default:
		http.Error(w, "Incorrect metric TYPE.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
