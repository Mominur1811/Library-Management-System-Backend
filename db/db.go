package db

import "github.com/jmoiron/sqlx"

var readDb *sqlx.DB
var writeDb *sqlx.DB

func GetReadDB() *sqlx.DB {
	return readDb
}

func GetWriteDB() *sqlx.DB {
	return writeDb
}
