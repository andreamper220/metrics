package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/server/handlers"
	"github.com/andreamper220/metrics.git/internal/server/middlewares"
	"github.com/andreamper220/metrics.git/internal/server/storages"
)

func MakeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Route(`/`, func(r chi.Router) {
		r.Get(`/`, middlewares.WithGzip(middlewares.WithLogging(handlers.ShowMetrics)))
		r.Post(`/value/`, middlewares.WithGzip(middlewares.WithLogging(handlers.ShowMetric)))
		r.Get(`/ping/`, middlewares.WithGzip(middlewares.WithLogging(handlers.Ping)))
	})
	r.Post(`/update/`, middlewares.WithGzip(middlewares.WithLogging(handlers.UpdateMetric)))

	// deprecated
	r.Get(`/value/{type}/{name}`, middlewares.WithGzip(middlewares.WithLogging(handlers.ShowMetricOld)))
	r.Post(`/update/{type}/{name}/{value}`, middlewares.WithGzip(middlewares.WithLogging(handlers.UpdateMetricOld)))

	return r
}

func MakeStorage() error {
	// choose metrics storage
	if Config.DatabaseDSN != "" {
		conn, err := sql.Open("pgx", Config.DatabaseDSN)
		if err == nil {
			storages.Storage = storages.NewDBStorage(conn)
		}
	} else if Config.FileStoragePath != "" {
		storages.Storage = storages.NewFileStorage(Config.FileStoragePath, Config.StoreInterval == 0)
		// to restore metrics from file
		if Config.Restore {
			if err := storages.Storage.ReadMetrics(); err != nil {
				return err
			}
		}
		// to store metrics to file
		blockDone := make(chan bool)
		if Config.StoreInterval > 0 {
			storeTicker := time.NewTicker(time.Duration(Config.StoreInterval) * time.Second)
			go func() {
				for {
					select {
					case <-storeTicker.C:
						if err := storages.Storage.WriteMetrics(); err != nil {
							logger.Log.Error(err.Error())
						}
					case <-blockDone:
						storeTicker.Stop()
						return
					}
				}
			}()
		}
	} else {
		storages.Storage = storages.NewMemStorage()
	}

	return nil
}

func Run() error {
	if err := logger.Initialize(); err != nil {
		return err
	}
	if err := MakeStorage(); err != nil {
		return err
	}

	return http.ListenAndServe(Config.ServerAddress.String(), MakeRouter())
}
