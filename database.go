package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Env struct {
	DB *sql.DB
}

const (
	host		= "db"
	port		= 5432
	user		= "avito-user"
	password	= "avito123"
	dbname		= "avito-test-task"
)

func DBConnect() (*sql.DB, error){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db, nil
}