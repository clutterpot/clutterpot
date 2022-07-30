package db

import (
	"log"
	"os"

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

	return db
}
