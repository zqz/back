package controllers

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo"
	"github.com/zqzca/back/models"
)

var GetParam = (*echo.Context).Param

func StartTransaction() *sql.Tx {
	db, err := models.GetDB()
	tx, err := db.Begin()

	if err != nil {
		fmt.Println("Failed to create transaction", err)
		return nil
	}

	return tx
}
