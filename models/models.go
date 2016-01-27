package models

import (
	"fmt"

	"upper.io/db"
)

var database db.Database

func SetDB(dbcopy *db.Database) {
	database = *dbcopy
	fmt.Println("thanks for the yummy db b0ss")
}
