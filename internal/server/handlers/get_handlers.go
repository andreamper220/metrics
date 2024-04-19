package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/go-chi/chi/v5"

	"github.com/andreamper220/metrics.git/internal/shared"
)

func ShowMetrics(w http.ResponseWriter, r *http.Request) {
	body := "== COUNTERS:\r\n"
	counterNames := make([]string, 0, len(storage.Counters))
	for name := range storage.Counters {
		counterNames = append(counterNames, string(name))
	}
	sort.Strings(counterNames)
	for _, name := range counterNames {
		body += fmt.Sprintf("= %s => %v\r\n", name, storage.Counters[shared.CounterMetricName(name)])
	}

	body += "== GAUGES:\r\n"
	gaugeNames := make([]string, 0, len(storage.Gauges))
	for name := range storage.Gauges {
		gaugeNames = append(gaugeNames, string(name))
	}
	sort.Strings(gaugeNames)
	for _, name := range gaugeNames {
		body += fmt.Sprintf("= %s => %v\r\n", name, storage.Gauges[shared.GaugeMetricName(name)])
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
	case shared.CounterMetricType:
		counterValue, ok := storage.Counters[shared.CounterMetricName(name)]
		if !ok {
			http.Error(w, "Incorrect metric NAME.", http.StatusNotFound)
			return
		}

		value = fmt.Sprintf("%d", counterValue)
	case shared.GaugeMetricType:
		gaugeValue, ok := storage.Gauges[shared.GaugeMetricName(name)]
		if !ok {
			http.Error(w, "Incorrect metric NAME.", http.StatusNotFound)
			return
		}

		value = fmt.Sprintf("%g", gaugeValue)
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