package lib

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func Connect() (*sqlx.DB, error) {
	open := os.Getenv("DATABASE_URL")

	if parsedURL, err := pq.ParseURL(open); err == nil && parsedURL != "" {
		open = parsedURL
	}

	con, err := sqlx.Connect("postgres", open)

	if err != nil {
		fmt.Println(err)
	}

	return con, err
}
