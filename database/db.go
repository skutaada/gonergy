package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	db, err := sql.Open("sqlite3", "/tmp/gonergy/db.sqlite")
	if err != nil {
		return err
	}
	
	sqlStmt := `
	create table if not exists energy (
		id integer not null primary key,
		date text not null unique,
		value real not null
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	DB = db
	return nil
}