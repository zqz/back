package db

import (
	"database/sql"
	"log"

	"github.com/DavidHuie/gomigrate"

	"upper.io/db"
	"upper.io/db/postgresql"
)

// AddDatabase connects to a psql db with the given name.
func DatabaseConnect(name string, user string) *db.Database {
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

	Migrate(database)

	// models.SetDB(&database)
	return &database
}

func Migrate(database db.Database) {
	drv := (database).Driver()

	m, err := gomigrate.NewMigrator(
		drv.(*sql.DB),
		gomigrate.Postgres{},
		"./migrations",
	)

	if err != nil {
		log.Fatalln("migrations failed", err.Error())
	}

	if err := m.Migrate(); err != nil {
		log.Fatalln("migrations failed", err.Error())
	}
}
