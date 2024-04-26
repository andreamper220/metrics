package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/server/handlers"
	"github.com/andreamper220/metrics.git/internal/server/middlewares"
)

func MakeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Route(`/`, func(r chi.Router) {
		r.Get(`/`, middlewares.WithGzip(middlewares.WithLogging(handlers.ShowMetrics)))
		r.Post(`/value/`, middlewares.WithGzip(middlewares.WithLogging(handlers.ShowMetric)))
	})
	r.Post(`/update/`, middlewares.WithGzip(middlewares.WithLogging(handlers.UpdateMetric)))

	// deprecated
	r.Get(`/value/{type}/{name}`, middlewares.WithGzip(middlewares.WithLogging(handlers.ShowMetricOld)))
	r.Post(`/update/{type}/{name}/{value}`, middlewares.WithGzip(middlewares.WithLogging(handlers.UpdateMetricOld)))

	return r
}

func Run() error {
	if err := logger.Initialize(); err != nil {
		return err
	}

	return http.ListenAndServe(Config.ServerAddress.String(), MakeRouter())
}
