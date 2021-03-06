package database

import (
	"database/sql"
	"log"
	"os"
	"time"
)

// DbConnection is the means by which we read/write data in our SQL database.
var DbConnection *sql.DB

// SetupDatabase will set up the connection to our SQL database.
func SetupDatabase() {
	var err error
	DbConnection, err = sql.Open("mysql", os.Getenv("MYSQL_URI"))
	if err != nil {
		log.Fatal(err)
	}
	DbConnection.SetMaxOpenConns(4)
	DbConnection.SetMaxIdleConns(4)
	DbConnection.SetConnMaxLifetime(60 * time.Second)
}
