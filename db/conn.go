package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DbConnection struct {
	db *sql.DB
}

func GetDbConnection() (*DbConnection, error) {
	db, err := sql.Open("sqlite3", "./lyrics.db")
	if err != nil {
		return nil, err
	}

	//ping db to make sure it's alive
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = prepDb(db)
	if err != nil {
		return nil, err
	}

	return &DbConnection{
		db: db,
	}, nil
}

func prepDb(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS lyrics (id INTEGER PRIMARY KEY, artist TEXT, title TEXT, lyrics TEXT, url TEXT)")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS queries (id INTEGER PRIMARY KEY, query TEXT, artist TEXT, title TEXT, url TEXT)")
	if err != nil {
		return err
	}

	return nil
}
