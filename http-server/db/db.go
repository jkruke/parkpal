package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("mysql", "guest:2107@tcp(192.168.0.198:3306)/iot_project")
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
