package lib

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/zqzca/back/db"
)

func Connect() error {
	open := os.Getenv("DATABASE_URL")

	if parsedURL, err := pq.ParseURL(open); err == nil && parsedURL != "" {
		open = parsedURL
	}

	con, err := sqlx.Connect("postgres", open)

	if err != nil {
		fmt.Println(err)
	}

	db.Connection = con

	return err
}
