package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Database struct {
	connection *sqlx.DB
}

func (v *Database) newConnection() {
	conn, err := sqlx.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	v.connection = conn
}

func (v *Database) Connection() *sqlx.DB {
	if v.connection == nil {
		v.newConnection()
	}
	return v.connection
}

func NewSqlLiteDatabase() *Database {
	return &Database{}
}
