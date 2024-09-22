package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const file string = "activities.db"

var conn *sql.DB

func InitDB() (*sql.DB, error) {
	if conn != nil {
		return conn, nil
	}
	var err error
	conn, err = sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return conn, nil
}

func InitTables() error {
	conn, err := InitDB()
	if err != nil {
		return err
	}
	const create string = `
	CREATE TABLE IF NOT EXISTS config (
	id INTEGER NOT NULL PRIMARY KEY,
	token TEXT NOT NULL
	);`

	if _, err := conn.Exec(create); err != nil {
		return err
	}
	return nil
}

func CloseDB() {
	if conn != nil {
		conn.Close()
		conn = nil
	}
}
