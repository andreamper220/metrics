package metric

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sort"
	"strconv"

	"github.com/andreamper220/metrics.git/internal/constants"
)

type MemStorage struct {
	counters map[string]int64
	gauges   map[string]float64
}

var storage = MemStorage{
	make(map[string]int64),
	make(map[string]float64),
}

type metric interface {
	store()
}

type counterMetric struct {
	name  string
	value int64
}

func (m *counterMetric) store() {
	storage.counters[m.name] += m.value
}

type gaugeMetric struct {
	name  string
	value float64
}

func (m *gaugeMetric) store() {
	storage.gauges[m.name] = m.value
}

func UpdateMetric(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "Not found metric NAME.", http.StatusNotFound)
		return
	}

	switch chi.URLParam(r, "type") {
	case constants.CounterMetricType:
		value, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
		if err != nil {
			http.Error(w, "Incorrect metric VALUE: "+err.Error(), http.StatusBadRequest)
			return
		}

		metric := counterMetric{name, value}
		metric.store()
	case constants.GaugeMetricType:
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

func ShowMetrics(w http.ResponseWriter, r *http.Request) {
	body := "== COUNTERS:\r\n"
	counterNames := make([]string, 0, len(storage.counters))
	for name := range storage.counters {
		counterNames = append(counterNames, name)
	}
	sort.Strings(counterNames)
	for _, name := range counterNames {
		body += fmt.Sprintf("= %s => %v\r\n", name, storage.counters[name])
	}

	body += "== GAUGES:\r\n"
	gaugeNames := make([]string, 0, len(storage.gauges))
	for name := range storage.gauges {
		gaugeNames = append(gaugeNames, name)
	}
	sort.Strings(gaugeNames)
	for _, name := range gaugeNames {
		body += fmt.Sprintf("= %s => %v\r\n", name, storage.gauges[name])
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ShowMetric(w http.ResponseWriter, r *http.Request) {
	var value string
	name := chi.URLParam(r, "name")

	switch chi.URLParam(r, "type") {
	case constants.CounterMetricType:
		counterValue, ok := storage.counters[name]
		if !ok {
			http.Error(w, "Incorrect metric NAME.", http.StatusNotFound)
			return
		}

		value = fmt.Sprintf("%d", counterValue)
	case constants.GaugeMetricType:
		gaugeValue, ok := storage.gauges[name]
		if !ok {
			http.Error(w, "Incorrect metric NAME.", http.StatusNotFound)
			return
		}

		value = fmt.Sprintf("%f", gaugeValue)
	default:
		http.Error(w, "Incorrect metric TYPE.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
