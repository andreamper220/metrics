package storages

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/andreamper220/metrics.git/internal/shared"
)

type DBStorage struct {
	*AbstractStorage
	Connection *sql.DB
}

func NewDBStorage(conn *sql.DB) (*DBStorage, error) {
	ctx := context.Background()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS metrics_counter (
		    id varchar(128) PRIMARY KEY NOT NULL,
		    value int NOT NULL
		)
	`)
	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS metrics_gauge (
		    id varchar(128) PRIMARY KEY NOT NULL,
		    value double precision NOT NULL
		)
	`)

	return &DBStorage{
		AbstractStorage: NewAbstractStorage(true),
		Connection:      conn,
	}, tx.Commit()
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

		err := retry.Do(
			func() error {
				_, err := conn.ExecContext(ctx, sqlString, sqlVars...)
				if err != nil {
					var pgErr *pgconn.PgError
					if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
						return err
					} else {
						return retry.Unrecoverable(err)
					}
				}
				return nil
			},
			retry.Attempts(3),
			retry.Delay(time.Second),
			retry.DelayType(retry.BackOffDelay),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
