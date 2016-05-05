package models

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

func Connection() (*sql.DB, error) {
	open := "dbname=zqz-test sslmode=disable"

	if parsedURL, err := pq.ParseURL(open); err == nil && parsedURL != "" {
		open = parsedURL
	}

	db, err := sql.Open("postgres", open)

	return db, err
}

// TxWrapper is used in tests
func TxWrapper(callback func(*sql.Tx)) {
	db, err := Connection()

	if err != nil {
		fmt.Println(err)
		return
	}

	if tx, err := db.Begin(); err == nil {
		callback(tx)
		tx.Rollback()
	} else {
		fmt.Println(err)
	}
}
