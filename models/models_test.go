package models

import (
	"database/sql"
	"fmt"

	"github.com/zqzca/back/db"
)

func init() {
	db := db.DatabaseConnect()
	SetDB(db)
}

func buildValidUser() *User {
	return &User{
		FirstName: "John",
		LastName:  "Carmack",
		Address:   "somewhere in texas",
		Username:  "jc",
		Phone:     "+123 123 1234",
		Email:     "johnc@idsoftware.com",
	}
}

func truncate(table string) {
	cmd := fmt.Sprintf("truncate %s;", table)
	database.Driver().(*sql.DB).Query(cmd)
}
