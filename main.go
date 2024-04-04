package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	username := "user"
	password := "password"
	dbname := "sources"
	dbHost := "sources_db"
	dbPort := 3306
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, dbHost, dbPort, dbname)

	_, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error connecting to the database %v", err)
	}
}
