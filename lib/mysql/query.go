package mysql

import (
	dbsql "database/sql"
	"hcc/flute/lib/config"
	"hcc/flute/lib/logger"
	"strings"
	"time"
)

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
func Prepare(sql string) (*dbsql.Stmt, error) {
	var stmt *dbsql.Stmt
	var err error

	for i := 0; i < int(config.Mysql.ConnectionRetryCount); i++ {
		stmt, err = Db.Prepare(sql)
		if err != nil {
			errMsg := strings.ToLower(err.Error())
			containInvalidConnection := strings.Contains(errMsg, "invalid connection")
			containBadConnection := strings.Contains(errMsg, "bad connection")

			if !containInvalidConnection && !containBadConnection {
				return nil, err
			}

			logger.Logger.Println("mysql.Prepare(): Retrying MySQL connection...")
		} else {
			return stmt, nil
		}

		time.Sleep(time.Millisecond * time.Duration(config.Mysql.ConnectionRetryIntervalMs))
	}

	return stmt, nil
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func Query(sql string, args ...interface{}) (*dbsql.Rows, error) {
	var rows *dbsql.Rows
	var err error

	for i := 0; i < int(config.Mysql.ConnectionRetryCount); i++ {
		rows, err = Db.Query(sql, args...)
		if err != nil {
			errMsg := strings.ToLower(err.Error())
			containInvalidConnection := strings.Contains(errMsg, "invalid connection")
			containBadConnection := strings.Contains(errMsg, "bad connection")

			if !containInvalidConnection && !containBadConnection {
				return nil, err
			}

			logger.Logger.Println("mysql.Query(): Retrying MySQL connection...")
		} else {
			return rows, nil
		}

		time.Sleep(time.Millisecond * time.Duration(config.Mysql.ConnectionRetryIntervalMs))
	}

	return rows, nil
}

// QueryRowScan copies the columns from the matched row into the values
// pointed at by dest. See the documentation on Rows.Scan for details.
// If more than one row matches the query,
// Scan uses the first row and discards the rest. If no row matches
// the query, Scan returns ErrNoRows.
func QueryRowScan(row *dbsql.Row, dest ...interface{}) error {
	var err error

	for i := 0; i < int(config.Mysql.ConnectionRetryCount); i++ {
		err = row.Scan(dest...)
		if err != nil {
			errMsg := strings.ToLower(err.Error())
			containInvalidConnection := strings.Contains(errMsg, "invalid connection")
			containBadConnection := strings.Contains(errMsg, "bad connection")

			if !containInvalidConnection && !containBadConnection {
				return err
			}

			logger.Logger.Println("mysql.QueryRowScan(): Retrying MySQL connection...")
		} else {
			return nil
		}

		time.Sleep(time.Millisecond * time.Duration(config.Mysql.ConnectionRetryIntervalMs))
	}

	return nil
}
