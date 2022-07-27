package db

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URI"))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Successfully connected to the database")

	if err := migrateDB(db); err != nil {
		log.Fatal(err)
	}

	return db
}

func migrateDB(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		log.Print("Migrations: ", err)
	}

	return nil
}
