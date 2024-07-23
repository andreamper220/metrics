package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

// ShowMetrics показывает все сохранённые ключ-значения метрик.
func ShowMetrics(w http.ResponseWriter, r *http.Request) {
	body := "== COUNTERS:\r\n"
	counters, err := storages.Storage.GetCounters()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for _, counter := range counters {
		body += fmt.Sprintf("= %s => %v\r\n", counter.Name, counter.Value)
	}

	body += "== GAUGES:\r\n"
	gauges, err := storages.Storage.GetGauges()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for _, gauge := range gauges {
		body += fmt.Sprintf("= %s => %v\r\n", gauge.Name, gauge.Value)
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ShowMetric отдаёт JSON с типом и значением искомой метрики.
func ShowMetric(w http.ResponseWriter, r *http.Request) {
	var metric shared.Metric

	// json decoder
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case shared.CounterMetricType:
		isExisted := false
		counters, err := storages.Storage.GetCounters()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for _, counter := range counters {
			if counter.Name == shared.CounterMetricName(metric.ID) {
				metric.Delta = &counter.Value
				isExisted = true
				break
			}
		}
		if !isExisted {
			http.Error(w, "Incorrect metric ID.", http.StatusNotFound)
			return
		}
	case shared.GaugeMetricType:
		isExisted := false
		gauges, err := storages.Storage.GetGauges()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for _, gauge := range gauges {
			if gauge.Name == shared.GaugeMetricName(metric.ID) {
				metric.Value = &gauge.Value
				isExisted = true
				break
			}
		}
		if !isExisted {
			http.Error(w, "Incorrect metric ID.", http.StatusNotFound)
			return
		}
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

// ShowMetricOld отдаёт значение искомой метрики.
func ShowMetricOld(w http.ResponseWriter, r *http.Request) {
	var value string
	name := chi.URLParam(r, "name")

	switch chi.URLParam(r, "type") {
	case shared.CounterMetricType:
		isExisted := false
		counters, err := storages.Storage.GetCounters()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for _, counter := range counters {
			if counter.Name == shared.CounterMetricName(name) {
				value = fmt.Sprintf("%d", counter.Value)
				isExisted = true
				break
			}
		}
		if !isExisted {
			http.Error(w, "Incorrect metric NAME.", http.StatusNotFound)
			return
		}
	case shared.GaugeMetricType:
		isExisted := false
		gauges, err := storages.Storage.GetGauges()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for _, gauge := range gauges {
			if gauge.Name == shared.GaugeMetricName(name) {
				value = fmt.Sprintf("%g", gauge.Value)
				isExisted = true
				break
			}
		}
		if !isExisted {
			http.Error(w, "Incorrect metric NAME.", http.StatusNotFound)
			return
		}
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

// Ping отдаёт 200.
func Ping(w http.ResponseWriter, r *http.Request) {
	storage, ok := storages.Storage.(*storages.DBStorage)
	if !ok {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
		return
	}
	if err := storage.Connection.Ping(); err != nil {
		http.Error(w, "DB storage ping error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
