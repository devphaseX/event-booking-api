package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	_DB, err := sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("db connection failed to start")
	}

	DB = _DB

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	createTables()
}

func createTables() {
	createEventTable := `
		CREATE TABLE IF NOT EXISTS events (
			id 	INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARHCAR(150) NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			dateTime DATETIME NOT NULL,
			user_id INTEGER
		)
	`

	_, err := DB.Exec(createEventTable)

	if err != nil {
		fmt.Println(err)
		panic("could not create event table.")
	}
}
