package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

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
	var metric shared.Metric

	// json decoder
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case shared.CounterMetricType:
		delta, ok := storage.Counters[shared.CounterMetricName(metric.ID)]
		if !ok {
			http.Error(w, "Incorrect metric ID.", http.StatusNotFound)
			return
		}

		metric.Delta = &delta
	case shared.GaugeMetricType:
		value, ok := storage.Gauges[shared.GaugeMetricName(metric.ID)]
		if !ok {
			http.Error(w, "Incorrect metric ID.", http.StatusNotFound)
			return
		}

		metric.Value = &value
	default:
		http.Error(w, "Incorrect metric TYPE.", http.StatusBadRequest)
		return
	}

	// json encoder
	if err := json.NewEncoder(w).Encode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
