package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)


func Connect() *sql.DB {
	connStr := "host=127.0.0.1 port=5434 user=user password=thisneedstochange dbname=webhookenginedb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	if err != nil {
		log.Fatal("here is the error ", err)
	}
	if err = db.Ping(); err != nil{
		log.Fatal("here is the error ", err)	
	}

	log.Println("Database connected successfully")
	return db
}