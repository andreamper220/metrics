package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/andreamper220/metrics.git/internal/server/handlers"
)

func MakeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Route(`/`, func(r chi.Router) {
		r.Get(`/`, handlers.ShowMetrics)
		r.Get(`/value/{type}/{name}`, handlers.ShowMetric)
	})
	r.Post(`/update/{type}/{name}/{value}`, handlers.UpdateMetric)

	return r
}

func Run() error {
	return http.ListenAndServe(Config.ServerAddress.String(), MakeRouter())
}
