package config

import (
	"database/sql"
)

var db *sql.DB

func GetDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	var err error
	db, err = sql.Open("sqlite3", "./thoth.db")
	if err != nil {
		return nil, err
	}

	println("Connected to database")

	return db, nil
}
