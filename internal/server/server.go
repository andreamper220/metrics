package server

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"

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
	})
	r.Post(`/update/`, middlewares.WithGzip(middlewares.WithLogging(handlers.UpdateMetric)))
	r.Post(`/updates/`, middlewares.WithGzip(middlewares.WithLogging(handlers.UpdateMetrics)))
	r.Get(`/ping`, middlewares.WithGzip(middlewares.WithLogging(handlers.Ping)))

	// deprecated
	r.Get(`/value/{type}/{name}`, middlewares.WithGzip(middlewares.WithLogging(handlers.ShowMetricOld)))
	r.Post(`/update/{type}/{name}/{value}`, middlewares.WithGzip(middlewares.WithLogging(handlers.UpdateMetricOld)))

	return r
}

func MakeStorage(blockDone chan bool) error {
	// choose metrics storage
	if Config.DatabaseDSN != "" {
		conn, err := sql.Open("pgx", Config.DatabaseDSN)
		if err == nil {
			storages.Storage, err = storages.NewDBStorage(conn)
			if err != nil {
				return err
			}
		}
	} else if Config.FileStoragePath != "" {
		var err error
		storages.Storage, err = storages.NewFileStorage(
			Config.FileStoragePath,
			Config.StoreInterval,
			Config.Restore,
			blockDone,
		)
		if err != nil {
			return err
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
	blockDone := make(chan bool)
	if err := MakeStorage(blockDone); err != nil {
		return err
	}
	if Config.DatabaseDSN != "" {
		storage, ok := storages.Storage.(*storages.DBStorage)
		if !ok {
			return errors.New("DB storage not created")
		}
		defer storage.Connection.Close()
	}

	err := http.ListenAndServe(Config.ServerAddress.String(), MakeRouter())
	<-blockDone
	return err
}
