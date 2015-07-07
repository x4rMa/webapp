package qmysql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type TMySQL struct {
	driver *sql.DB
}

func New(driver, source string) (db TMySQL, err error) {

	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%s", e))
		}
	}()

	db.driver, err = sql.Open(driver, source)
	if err != nil {
		panic("Open(): " + err.Error())
	}

	err = db.driver.Ping()
	if err != nil {
		panic("Ping(): " + err.Error())
	}

	return db, nil
}

func sqlFetch(Rows *sql.Rows) ([]map[string]string, error) {

	columns, err := Rows.Columns()
	if err != nil {
		return nil, err
	}

	if len(columns) == 0 {
		return nil, errors.New("SQLFetch(): get columns error")
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	data := make([]map[string]string, 0)

	for Rows.Next() {
		newRow := make(map[string]string)

		if err := Rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			newRow[columns[i]] = value
		}
		data = append(data, newRow)
	}

	return data, nil
}

func (db TMySQL) Query(SQL string, args ...interface{}) ([]map[string]string, error) {
	var (
		rows *sql.Rows
		err  error
	)

	rows, err = db.driver.Query(SQL, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return sqlFetch(rows)
}

func (db TMySQL) Result(sql string, args ...interface{}) (string, error) {
	res, err := db.Query(sql, args)
	if err != nil {
		return "", err
	}

	for _, val := range res[0] {
		return val, nil
	}

	return "NEVER", nil
}

func (db TMySQL) Exec(SQL string, args ...interface{}) (interface{}, error) {
	var (
		res sql.Result
		err error
	)

	res, err = db.driver.Exec(SQL, args...)

	return res, err
}

func (db TMySQL) ExecId(SQL string, args ...interface{}) (int64, error) {
	var (
		res sql.Result
		err error
	)

	res, err = db.driver.Exec(SQL, args...)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (db TMySQL) Start() {
	db.Query("BEGIN")
}

func (db TMySQL) Rollback() {
	db.Query("ROLLBACK")
}

func (db TMySQL) Commit() {
	db.Query("COMMIT")
}
