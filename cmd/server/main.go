package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	metric "github.com/andreamper220/metrics.git/internal/server"
)

func main() {
	r := chi.NewRouter()
	r.Route(`/`, func(r chi.Router) {
		r.Get(`/`, metric.ShowMetrics)
		r.Get(`/value/{type}/{name}`, metric.ShowMetric)
	})
	r.Post(`/update/{type}/{name}/{value}`, metric.UpdateMetric)

	log.Fatal(http.ListenAndServe(":8080", r))
}
