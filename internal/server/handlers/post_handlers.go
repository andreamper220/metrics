package handlers

import (
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

		metric := counterMetric{name, value}
		metric.store()
	case shared.GaugeMetricType:
		value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
		if err != nil {
			http.Error(w, "Incorrect metric VALUE: "+err.Error(), http.StatusBadRequest)
			return
		}

		metric := gaugeMetric{name, value}
		metric.store()
	default:
		http.Error(w, "Incorrect metric TYPE.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
