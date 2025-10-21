package sql_source

import (
	"database/sql"
	"log"
	"project/internals/data/config"

	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	env, err := config.NewEnv()
	if err != nil {

		return nil
	}
	er := env.GetValueForKey("DB_PASS")
	if er == "" {
		return nil
	}
	db, err := sql.Open("postgres", "postgres://postgres:pass@127.0.0.1:5432/docsniff_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
