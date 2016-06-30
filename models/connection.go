package models

// import (
// 	"database/sql"
// 	"fmt"
// 	"os"
// 	"sync"

// 	"github.com/lib/pq"
// )

// type Executor interface {
// 	Exec(query string, args ...interface{}) (sql.Result, error)
// 	QueryRow(query string, args ...interface{}) *sql.Row
// 	Query(query string, args ...interface{}) (*sql.Rows, error)
// }

// var database *sql.DB
// var mu sync.Mutex

// func GetDB() (*sql.DB, error) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	var err error
// 	if database == nil {
// 		database, err = Connection()
// 	}

// 	return database, err
// }

// func Connection() (*sql.DB, error) {
// 	open := os.Getenv("DATABASE_URL")

// 	if parsedURL, err := pq.ParseURL(open); err == nil && parsedURL != "" {
// 		open = parsedURL
// 	}

// 	db, err := sql.Open("postgres", open)

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	return db, err
// }

// // TxWrapper is used in tests
// func TxWrapper(callback func(*sql.Tx)) {
// 	db, err := GetDB()

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	if tx, err := db.Begin(); err == nil {

// 		defer func() {
// 			// Catch panics
// 			if r := recover(); r != nil {
// 				tx.Rollback()
// 			}
// 		}()

// 		callback(tx)
// 		tx.Rollback()
// 	} else {
// 		fmt.Println(err)
// 	}
// }
