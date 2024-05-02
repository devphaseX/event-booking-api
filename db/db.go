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
	createUserTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email VARCHAR(150) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			passwordSalt TEXT NOT NULL
		)
	`

	_, err := DB.Exec(createUserTable)

	if err != nil {
		fmt.Println(err)
		panic("could not create user table.")
	}

	createEventTable := `
		CREATE TABLE IF NOT EXISTS events (
			id 	INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARHCAR(150) NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			dateTime DATETIME NOT NULL,
			user_id INTEGER REFERENCES users(id)
	
		)
	`
	_, err = DB.Exec(createEventTable)

	if err != nil {
		fmt.Println(err)
		panic("could not create event table.")
	}

	createTicketTable := `
		CREATE TABLE IF NOT EXISTS tickets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER REFERENCES users(id),
			event_id INTEGER,
			FOREIGN KEY(event_id) REFERENCES events(id)
		)
	`

	_, err = DB.Exec(createTicketTable)

	if err != nil {
		fmt.Println(err)
		panic("could not create ticket table.")
	}
}
