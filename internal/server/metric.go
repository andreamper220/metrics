package metric

import (
	"fmt"
	"net/http"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "Not found metric NAME.", http.StatusNotFound)
		return
	}

	switch r.PathValue("type") {
	case constants.CounterMetricType:
		value, err := strconv.ParseInt(r.PathValue("value"), 10, 64)
		if err != nil {
			http.Error(w, "Incorrect metric VALUE: "+err.Error(), http.StatusBadRequest)
			return
		}

		metric := counterMetric{name, value}
		metric.store()
		fmt.Printf("[%s => %v] metric is in storage\n", name, storage.counters[name])
	case constants.GaugeMetricType:
		value, err := strconv.ParseFloat(r.PathValue("value"), 64)
		if err != nil {
			http.Error(w, "Incorrect metric VALUE: "+err.Error(), http.StatusBadRequest)
			return
		}

		metric := gaugeMetric{name, value}
		metric.store()
		fmt.Printf("[%s => %v] metric is in storage\n", name, storage.gauges[name])
	default:
		http.Error(w, "Incorrect metric TYPE.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
