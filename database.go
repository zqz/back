package main

import (
	"log"

	"github.com/zqzca/back/models"

	"upper.io/db"
	"upper.io/db/postgresql"
)

var database db.Database

// AddDatabase connects to a psql db with the given name.
func DatabaseConnect(name string, user string) {
	settings := postgresql.ConnectionURL{
		Database: name,
		User:     user,
	}

	var err error
	database, err := db.Open(postgresql.Adapter, settings)

	if err != nil {
		log.Fatalf("Failed to connect to database: %s with user %s - %s\n", name, user, err.Error())
	}

	if err = database.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %s - %s\n", name, err.Error())
	}

	models.SetDB(&database)
}
