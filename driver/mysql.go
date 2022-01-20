package driver

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func ConnectToSQL() (*sql.DB, error) {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "organisation",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected!")

	return db, nil
}
