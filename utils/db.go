package utils

import (
	"database/sql"
	"fmt"

	//pq is the sql driver for database/sql
	_ "github.com/lib/pq"
)

const (
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "postgres"
	dbPort     = 5432
)

// Connect - connects to DB
func Connect() (*sql.DB, error) {

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable", dbUser, dbPassword, dbName, dbPort)
	db, err := sql.Open("postgres", dbinfo)
	return db, err
}

// CloseConnection - closes DB connection
func CloseConnection(db *sql.DB) {
	db.Close()
}
