package models_test

import (
	"github.com/zqzca/back/db"
	. "github.com/zqzca/back/models"
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

func buildFile() *File {
	return &File{
		Name:  "Foo",
		Size:  123,
		Hash:  "foo",
		Type:  "image/jpg",
		State: Incomplete,
	}
}

func createFile() *File {
	f := buildFile()
	f.Save()
	return f
}
