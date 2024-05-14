package storages

import (
	"context"
	"database/sql"
	"embed"
	"strconv"

	"github.com/pressly/goose/v3"

	"github.com/andreamper220/metrics.git/internal/shared"
)

//go:embed migrations/*.sql
var migrations embed.FS

type DBStorage struct {
	metrics            metrics
	toSaveMetricsAsync bool
	Connection         *sql.DB
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
		metrics: metrics{
			counters: make(map[shared.CounterMetricName]int64),
			gauges:   make(map[shared.GaugeMetricName]float64),
		},
		toSaveMetricsAsync: true,
		Connection:         conn,
	}, nil
}
func (dbs *DBStorage) GetCounters() map[shared.CounterMetricName]int64 {
	return dbs.metrics.counters
}
func (dbs *DBStorage) SetCounters(counters map[shared.CounterMetricName]int64) error {
	for name, value := range counters {
		dbs.metrics.counters[name] = value
	}
	return nil
}
func (dbs *DBStorage) GetGauges() map[shared.GaugeMetricName]float64 {
	return dbs.metrics.gauges
}
func (dbs *DBStorage) SetGauges(gauges map[shared.GaugeMetricName]float64) error {
	for name, value := range gauges {
		dbs.metrics.gauges[name] = value
	}
	return nil
}
func (dbs *DBStorage) GetToSaveMetricsAsync() bool {
	return dbs.toSaveMetricsAsync
}
func (dbs *DBStorage) WriteMetrics() error {
	if err := insertMetrics(context.Background(), dbs.Connection, dbs.metrics.counters, "metrics_counter"); err != nil {
		return err
	}
	if err := insertMetrics(context.Background(), dbs.Connection, dbs.metrics.gauges, "metrics_gauge"); err != nil {
		return err
	}

	return nil
}
func (dbs *DBStorage) ReadMetrics() error {
	ctx := context.Background()

	counterRows, err := dbs.Connection.QueryContext(ctx, `
		SELECT id, value
		FROM metrics_counter
	`)
	if err != nil {
		return err
	}
	defer counterRows.Close()

	gaugeRows, err := dbs.Connection.QueryContext(ctx, `
		SELECT id, value
		FROM metrics_gauge
	`)
	if err != nil {
		return err
	}
	defer gaugeRows.Close()

	for counterRows.Next() {
		var id string
		var value int64
		if err := counterRows.Scan(&id, &value); err != nil {
			return err
		}
		dbs.metrics.counters[shared.CounterMetricName(id)] = value
	}
	for gaugeRows.Next() {
		var id string
		var value float64
		if err := gaugeRows.Scan(&id, &value); err != nil {
			return err
		}
		dbs.metrics.gauges[shared.GaugeMetricName(id)] = value
	}
	if err := counterRows.Err(); err != nil {
		return err
	}
	if err = gaugeRows.Err(); err != nil {
		return err
	}

	return nil
}

func insertMetrics[K shared.CounterMetricName | shared.GaugeMetricName, V int64 | float64](
	ctx context.Context, conn *sql.DB, metrics map[K]V, tableName string,
) error {
	if len(metrics) > 0 {
		sqlVarNumber := 1
		sqlString := "INSERT INTO " + tableName + " (id, value) VALUES "
		sqlVars := make([]any, len(metrics)*2)
		for name, value := range metrics {
			sqlString += "($" + strconv.Itoa(sqlVarNumber) + ", $" + strconv.Itoa(sqlVarNumber+1) + "),"
			sqlVars[sqlVarNumber-1] = name
			sqlVars[sqlVarNumber] = value
			sqlVarNumber += 2
		}
		sqlString = sqlString[0 : len(sqlString)-1]
		sqlString += " ON CONFLICT (id) DO UPDATE SET value = excluded.value;"

		_, err := conn.ExecContext(ctx, sqlString, sqlVars...)
		if err != nil {
			return err
		}
	}

	return nil
}
