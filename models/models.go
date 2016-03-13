package models

import (
	"database/sql"
	"fmt"

	"upper.io/db"
)

var Database db.Database

func SetDB(dbcopy *db.Database) {
	Database = *dbcopy
	fmt.Println("thanks for the yummy db b0ss")
}

// Truncate is only used for testing.
func Truncate(table string) {
	cmd := fmt.Sprintf("truncate %s cascade;", table)
	fmt.Println(cmd)
	_, err := Database.Driver().(*sql.DB).Query(cmd)

	if err != nil {
		fmt.Println("error: %s", err)
	}
}
