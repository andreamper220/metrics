package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/go-chi/chi/v5"

	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

func ShowMetrics(w http.ResponseWriter, r *http.Request) {
	body := "== COUNTERS:\r\n"
	counterNames := make([]string, 0, len(storages.Storage.GetCounters()))
	for name := range storages.Storage.GetCounters() {
		counterNames = append(counterNames, string(name))
	}
	sort.Strings(counterNames)
	for _, name := range counterNames {
		body += fmt.Sprintf("= %s => %v\r\n", name, storages.Storage.GetCounters()[shared.CounterMetricName(name)])
	}

	body += "== GAUGES:\r\n"
	gaugeNames := make([]string, 0, len(storages.Storage.GetGauges()))
	for name := range storages.Storage.GetGauges() {
		gaugeNames = append(gaugeNames, string(name))
	}
	sort.Strings(gaugeNames)
	for _, name := range gaugeNames {
		body += fmt.Sprintf("= %s => %v\r\n", name, storages.Storage.GetGauges()[shared.GaugeMetricName(name)])
	}

	w.Header().Set("Content-Type", "text/html")
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
		delta, ok := storages.Storage.GetCounters()[shared.CounterMetricName(metric.ID)]
		if !ok {
			http.Error(w, "Incorrect metric ID.", http.StatusNotFound)
			return
		}

		metric.Delta = &delta
	case shared.GaugeMetricType:
		value, ok := storages.Storage.GetGauges()[shared.GaugeMetricName(metric.ID)]
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ShowMetricOld(w http.ResponseWriter, r *http.Request) {
	var value string
	name := chi.URLParam(r, "name")

	switch chi.URLParam(r, "type") {
	case shared.CounterMetricType:
		counterValue, ok := storages.Storage.GetCounters()[shared.CounterMetricName(name)]
		if !ok {
			http.Error(w, "Incorrect metric NAME.", http.StatusNotFound)
			return
		}

		value = fmt.Sprintf("%d", counterValue)
	case shared.GaugeMetricType:
		gaugeValue, ok := storages.Storage.GetGauges()[shared.GaugeMetricName(name)]
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

func Ping(w http.ResponseWriter, r *http.Request) {
	storage, ok := storages.Storage.(*storages.DBStorage)
	if !ok {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}
	if err := storage.Connection.Ping(); err != nil {
		http.Error(w, "DB storage not created", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
