package db

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS sim_run (
		id	INTEGER NOT NULL PRIMARY KEY,
		time	DATETIME,
	);
	CREATE TABLE IF NOT EXISTS 
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
}
