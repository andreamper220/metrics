package storages

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"

	"github.com/andreamper220/metrics.git/internal/shared"
)

//go:embed migrations/*.sql
var migrations embed.FS

type DBStorage struct {
	Connection *sql.DB
}

func NewDBStorage(conn *sql.DB) (*DBStorage, error) {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return nil, err
	}
	if err := goose.Up(conn, "migrations"); err != nil {
		return nil, err
	}

	return &DBStorage{
		Connection: conn,
	}, nil
}
func (dbs *DBStorage) GetCounters() ([]CounterMetric, error) {
	ctx := context.Background()
	metrics := make([]CounterMetric, 0)

	counterRows, err := dbs.Connection.QueryContext(ctx, `
		SELECT id, value
		FROM metrics_counter
	`)
	if err != nil {
		return nil, err
	}
	defer counterRows.Close()

	for counterRows.Next() {
		var id string
		var value int64
		if err := counterRows.Scan(&id, &value); err != nil {
			return nil, err
		}
		metrics = append(metrics, CounterMetric{
			Name:  shared.CounterMetricName(id),
			Value: value,
		})
	}
	if err := counterRows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}
func (dbs *DBStorage) AddCounter(metric CounterMetric) error {
	if err := insertMetric(
		context.Background(), dbs.Connection, metric.Name, metric.Value, "metrics_counter",
	); err != nil {
		return err
	}
	return nil
}
func (dbs *DBStorage) GetGauges() ([]GaugeMetric, error) {
	ctx := context.Background()
	metrics := make([]GaugeMetric, 0)

	gaugeRows, err := dbs.Connection.QueryContext(ctx, `
		SELECT id, value
		FROM metrics_gauge
	`)
	if err != nil {
		return nil, err
	}
	defer gaugeRows.Close()

	for gaugeRows.Next() {
		var id string
		var value float64
		if err := gaugeRows.Scan(&id, &value); err != nil {
			return nil, err
		}
		metrics = append(metrics, GaugeMetric{
			Name:  shared.GaugeMetricName(id),
			Value: value,
		})
	}
	if err := gaugeRows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}
func (dbs *DBStorage) AddGauge(metric GaugeMetric) error {
	if err := insertMetric(
		context.Background(), dbs.Connection, metric.Name, metric.Value, "metrics_gauge",
	); err != nil {
		return err
	}
	return nil
}
func (dbs *DBStorage) GetMetrics() (Metrics, error) {
	counterMetrics, err := dbs.GetCounters()
	if err != nil {
		return Metrics{}, err
	}
	gaugeMetrics, err := dbs.GetGauges()
	if err != nil {
		return Metrics{}, err
	}

	return Metrics{
		counters: counterMetrics,
		gauges:   gaugeMetrics,
	}, nil
}

func insertMetric[K shared.CounterMetricName | shared.GaugeMetricName, V int64 | float64](
	ctx context.Context, conn *sql.DB, name K, value V, tableName string,
) error {
	_, err := conn.ExecContext(
		ctx,
		"INSERT INTO "+tableName+" (id, value) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET value = excluded.value;",
		name, value,
	)
	if err != nil {
		return err
	}

	return nil
}

func insertMetrics(
	ctx context.Context, conn *sql.DB, counterMetrics []CounterMetric, gaugeMetrics []GaugeMetric,
) error {
	var err error
	if len(counterMetrics) > 0 {
		for _, metric := range counterMetrics {
			err = insertMetric(ctx, conn, metric.Name, metric.Value, "metrics_counter")
		}
		if err != nil {
			return err
		}
	}

	if len(gaugeMetrics) > 0 {
		for _, metric := range counterMetrics {
			err = insertMetric(ctx, conn, metric.Name, metric.Value, "metrics_gauges")
		}
		if err != nil {
			return err
		}
	}

	return nil
}
