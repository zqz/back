package controllers

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/zqzca/back/models"
)

var GetParam = (*echo.Context).Param
var database *sql.DB

func init() {
	var err error
	database, err = models.Connection()

	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		os.Exit(0)
	}
}

func StartTransaction() *sql.Tx {
	tx, err := database.Begin()

	if err != nil {
		fmt.Println("Failed to create transaction", err)
		return nil
	}

	return tx
}
