package server

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/andreamper220/metrics.git/internal/logger"
	"github.com/andreamper220/metrics.git/internal/server/handlers"
	"github.com/andreamper220/metrics.git/internal/server/middlewares"
	"github.com/andreamper220/metrics.git/internal/server/storages"
	"github.com/andreamper220/metrics.git/internal/shared"
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
	storages.FileStoragePath = Config.FileStoragePath
	if Config.StoreInterval == 0 {
		storages.ToSaveMetricsAsync = true
	}

	// to restore metrics from file
	if Config.Restore {
		file, err := os.OpenFile(Config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer func() {
			_ = file.Close()
		}()

		fr := bufio.NewReader(file)
		dec := json.NewDecoder(fr)
		for {
			var metric shared.Metric

			err := dec.Decode(&metric)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			switch metric.MType {
			case shared.CounterMetricType:
				storages.Storage.Counters[shared.CounterMetricName(metric.ID)] = *metric.Delta
			case shared.GaugeMetricType:
				storages.Storage.Gauges[shared.GaugeMetricName(metric.ID)] = *metric.Value
			default:
				logger.Log.Fatalf("incorrect metric: %s", metric.ID)
			}
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
					if err := storages.Storage.StoreMetrics(); err != nil {
						logger.Log.Error(err.Error())
					}
				case <-blockDone:
					storeTicker.Stop()
					return
				}
			}
		}()
	}

	err := http.ListenAndServe(Config.ServerAddress.String(), MakeRouter())
	if Config.StoreInterval > 0 {
		<-blockDone
	}
	return err
}
