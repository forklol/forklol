package db

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var conn *sqlx.DB

func InitDB(connstr string) {
	db, err := sqlx.Connect("mysql", connstr)
	if err != nil {
		log.Fatalf("%s / %s\n", err, connstr)
	}

	conn = db
}

func GetDB() *sqlx.DB {
	return conn
}
