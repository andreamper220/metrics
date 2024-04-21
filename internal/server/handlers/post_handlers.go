package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
	"net/http"
)

var storage = &storages.MemStorage{
	Counters: make(map[shared.CounterMetricName]int64),
	Gauges:   make(map[shared.GaugeMetricName]float64),
}

type counterMetric struct {
	name  string
	value int64
}

func (m *counterMetric) store() {
	storage.Counters[shared.CounterMetricName(m.name)] += m.value
}

type gaugeMetric struct {
	name  string
	value float64
}

func (m *gaugeMetric) store() {
	storage.Gauges[shared.GaugeMetricName(m.name)] = m.value
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
		metric := counterMetric{reqMetric.ID, *reqMetric.Delta}
		metric.store()
	case shared.GaugeMetricType:
		metric := gaugeMetric{reqMetric.ID, *reqMetric.Value}
		metric.store()
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
