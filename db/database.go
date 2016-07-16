package db

import (
	"database/sql"
	"fmt"
	"sync"
)

type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

var Connection *sql.DB
var mu sync.Mutex

func StartTransaction() *sql.Tx {
	tx, err := Connection.Begin()

	if err != nil {
		fmt.Println("Failed to create transaction", err)
		return nil
	}

	return tx
}

// TxWrapper is used in tests
func TxWrapper(callback func(Executor)) {
	if tx, err := Connection.Begin(); err == nil {
		defer func() {
			// Catch panics
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		callback(tx)
		tx.Rollback()
	} else {
		fmt.Println(err)
	}
}
