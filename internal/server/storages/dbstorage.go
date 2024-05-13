package storages

import (
	"context"
	"database/sql"
	"github.com/andreamper220/metrics.git/internal/shared"
	"strconv"
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
	ctx := context.Background()

	if len(dbs.metrics.counters) > 0 {
		sqlVarNumber := 1
		sqlString := "INSERT INTO metrics_counter (id, value) VALUES "
		sqlVars := make([]any, len(dbs.metrics.counters)*2)
		for name, value := range dbs.metrics.counters {
			sqlString += "($" + strconv.Itoa(sqlVarNumber) + ", $" + strconv.Itoa(sqlVarNumber+1) + "),"
			sqlVars[sqlVarNumber-1] = name
			sqlVars[sqlVarNumber] = value
			sqlVarNumber += 2
		}
		sqlString = sqlString[0 : len(sqlString)-1]
		sqlString += " ON CONFLICT (id) DO UPDATE SET value = excluded.value;"

		_, err := dbs.Connection.ExecContext(ctx, sqlString, sqlVars...)
		if err != nil {
			return err
		}
	}

	if len(dbs.metrics.gauges) > 0 {
		sqlVarNumber := 1
		sqlString := "INSERT INTO metrics_gauge (id, value) VALUES "
		sqlVars := make([]any, len(dbs.metrics.gauges)*2)
		for name, value := range dbs.metrics.gauges {
			sqlString += "($" + strconv.Itoa(sqlVarNumber) + ", $" + strconv.Itoa(sqlVarNumber+1) + "),"
			sqlVars[sqlVarNumber-1] = name
			sqlVars[sqlVarNumber] = value
			sqlVarNumber += 2
		}
		sqlString = sqlString[0 : len(sqlString)-1]
		sqlString += " ON CONFLICT (id) DO UPDATE SET value = excluded.value;"

		_, err := dbs.Connection.ExecContext(ctx, sqlString, sqlVars...)
		if err != nil {
			return err
		}
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
