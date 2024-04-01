package main

import (
	"log"
	"net/http"

	metric "github.com/andreamper220/metrics.git/internal/server"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/{type}/{name}/{value}`, metric.UpdateMetric)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
