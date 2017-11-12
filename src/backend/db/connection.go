package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var conn *sqlx.DB

func InitDB(connstr string) error {
	db, err := sqlx.Connect("mysql", connstr)
	if err != nil {
		return err
	}

	conn = db
	return nil
}

func GetDB() *sqlx.DB {
	return conn
}
