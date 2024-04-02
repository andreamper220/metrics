package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	config "github.com/andreamper220/metrics.git/internal/config/server"
	metric "github.com/andreamper220/metrics.git/internal/server"
)

func main() {
	config.ParseFlags()

	r := chi.NewRouter()
	r.Route(`/`, func(r chi.Router) {
		r.Get(`/`, metric.ShowMetrics)
		r.Get(`/value/{type}/{name}`, metric.ShowMetric)
	})
	r.Post(`/update/{type}/{name}/{value}`, metric.UpdateMetric)

	serverAddress := config.Config.ServerAddress.String()
	fmt.Println("Running server on", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}
