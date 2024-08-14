package application

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
	"net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/server/application/handlers"
	"github.com/andreamper220/metrics.git/internal/server/application/middlewares"
	"github.com/andreamper220/metrics.git/internal/server/infrastructure/storages"
)

func MakeRouter() *chi.Mux {
	r := chi.NewRouter()
	// "show" routes
	r.Route(`/`, func(r chi.Router) {
		r.Get(`/`, middlewares.WithGzip(middlewares.WithLogging(handlers.ShowMetrics)))
		r.Post(`/value/`, middlewares.WithGzip(middlewares.WithLogging(handlers.ShowMetric)))
	})
	// "update" routes
	updateMetric := middlewares.WithGzip(middlewares.WithLogging(handlers.UpdateMetric))
	updateMetrics := middlewares.WithGzip(middlewares.WithLogging(handlers.UpdateMetrics))
	if Config.Sha256Key != "" {
		updateMetric = middlewares.WithSha256(updateMetric, Config.Sha256Key)
		updateMetrics = middlewares.WithSha256(updateMetrics, Config.Sha256Key)
	}
	if Config.CryptoKeyPath != "" {
		updateMetric = middlewares.WithCrypto(updateMetric, Config.CryptoKeyPath)
		updateMetrics = middlewares.WithCrypto(updateMetric, Config.CryptoKeyPath)
	}
	r.Post(`/update/`, updateMetric)
	r.Post(`/updates/`, updateMetrics)
	r.Get(`/ping`, middlewares.WithGzip(middlewares.WithLogging(handlers.Ping)))

	// deprecated
	r.Get(`/value/{type}/{name}`, middlewares.WithGzip(middlewares.WithLogging(handlers.ShowMetricOld)))
	r.Post(`/update/{type}/{name}/{value}`, middlewares.WithGzip(middlewares.WithLogging(handlers.UpdateMetricOld)))

	// service
	r.Get(`/debug/pprof/heap`, pprof.Handler("heap").ServeHTTP)

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

func Run(serverless bool) error {
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

	if serverless {
		return nil
	}

	var srv = http.Server{Addr: Config.ServerAddress.String(), Handler: MakeRouter()}

	idleConnsClosed := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	go func() {
		<-ctx.Done()
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Log.Error("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
		close(blockDone)
	}()
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
	<-blockDone
	return nil
}
