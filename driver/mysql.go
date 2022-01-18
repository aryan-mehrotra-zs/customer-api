package driver

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectToSQL() *sql.DB {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "organisation",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Println(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	log.Println("Connected!")

	return db
}
