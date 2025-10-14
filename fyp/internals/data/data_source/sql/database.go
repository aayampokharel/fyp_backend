package sql_source

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:rakeshyadav@9898@127.0.0.1:5432/docsniff_db?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}
	return db
}
