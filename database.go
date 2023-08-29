package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Env struct {
	DB *sql.DB
}

func DBConnect() (*sql.DB, error){
	connStr := os.Getenv("POSTGRE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}