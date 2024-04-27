package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/andreamper220/metrics.git/internal/logger"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-chi/chi/v5"

	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
)

func ShowMetrics(w http.ResponseWriter, r *http.Request) {
	//body := "== COUNTERS:\r\n"
	counterNames := make([]string, 0, len(storages.Storage.Counters))
	for name := range storages.Storage.Counters {
		counterNames = append(counterNames, string(name))
	}
	sort.Strings(counterNames)
	counters := make(map[string]int64, 5)
	for _, name := range counterNames {
		counters[name] = storages.Storage.Counters[shared.CounterMetricName(name)]
		//body += fmt.Sprintf("= %s => %v\r\n", name, storages.Storage.Counters[shared.CounterMetricName(name)])
	}

	//body += "== GAUGES:\r\n"
	gaugeNames := make([]string, 0, len(storages.Storage.Gauges))
	for name := range storages.Storage.Gauges {
		gaugeNames = append(gaugeNames, string(name))
	}
	sort.Strings(gaugeNames)
	gauges := make(map[string]float64, 30)
	for _, name := range gaugeNames {
		gauges[name] = storages.Storage.Gauges[shared.GaugeMetricName(name)]
		//body += fmt.Sprintf("= %s => %v\r\n", name, storages.Storage.Gauges[shared.GaugeMetricName(name)])
	}

	type bodyStruct struct {
		Counters map[string]int64
		Gauges   map[string]float64
	}
	body := &bodyStruct{
		Counters: counters,
		Gauges:   gauges,
	}

	dir, _ := filepath.Split(os.Args[0])
	filePath := filepath.Join(dir, "internal/templates/show_metrics.html")
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, body)
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
		delta, ok := storages.Storage.Counters[shared.CounterMetricName(metric.ID)]
		if !ok {
			http.Error(w, "Incorrect metric ID.", http.StatusNotFound)
			return
		}

		metric.Delta = &delta
	case shared.GaugeMetricType:
		value, ok := storages.Storage.Gauges[shared.GaugeMetricName(metric.ID)]
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
		counterValue, ok := storages.Storage.Counters[shared.CounterMetricName(name)]
		if !ok {
			http.Error(w, "Incorrect metric NAME.", http.StatusNotFound)
			return
		}

		value = fmt.Sprintf("%d", counterValue)
	case shared.GaugeMetricType:
		gaugeValue, ok := storages.Storage.Gauges[shared.GaugeMetricName(name)]
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
