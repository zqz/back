package models

import "github.com/zqzca/back/db"

func init() {
	db := db.DatabaseConnect("zqz-users-test", "dylan")
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
