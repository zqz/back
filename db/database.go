package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/DavidHuie/gomigrate"

	"upper.io/db"
	"upper.io/db/postgresql"
)

// AddDatabase connects to a psql db with the given name.
func DatabaseConnect() *db.Database {
	name := os.Getenv("DATABASE_NAME")
	if name == "" {
		fmt.Println("DATABASE_NAME not specified")
		os.Exit(1)
	}

	user := os.Getenv("DATABASE_USER")
	if user == "" {
		fmt.Println("DATABASE_USER not specified")
		os.Exit(1)
	}

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

	migrationPath := os.Getenv("MIGRATIONS")

	if len(migrationPath) > 0 {
		Migrate(database, migrationPath)
	}

	// models.SetDB(&database)
	return &database
}

func Migrate(database db.Database, migrationsRoot string) {
	path := fmt.Sprintf("%s/%s", migrationsRoot, "migrations")

	if len(path) == 0 {
		log.Fatalln("No migration path specified")
	}

	drv := (database).Driver()

	m, err := gomigrate.NewMigrator(
		drv.(*sql.DB),
		gomigrate.Postgres{},
		path,
	)

	if err != nil {
		log.Fatalln("migrations failed", err.Error())
	}

	if err := m.Migrate(); err != nil {
		log.Fatalln("migrations failed", err.Error())
	}
}
