package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbMutex sync.Mutex

func CreateTables() error {

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS tournaments (
	    id        TEXT NOT NULL PRIMARY KEY,
	    state     INTEGER,
	    start_time DATETIME,
	    end_time  DATETIME
	);

	CREATE TABLE IF NOT EXISTS matchups (
	    id          TEXT NOT NULL PRIMARY KEY,
	    state       INTEGER,
	    white       TEXT,
	    black       TEXT,
	    rounds      INTEGER,
	    completed   INTEGER,
	    black_wins  INTEGER,
	    white_wins  INTEGER,
	    draws       INTEGER,
	    errors      INTEGER
	);

	CREATE TABLE IF NOT EXISTS games (
	    id        TEXT NOT NULL PRIMARY KEY,
	    state     INTEGER,
	    start_time DATETIME,
	    end_time  DATETIME,
	    moves     TEXT,
	    outcome   TEXT,
	    method    TEXT
	);

	CREATE TABLE IF NOT EXISTS tournament_matchup (
	    id            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	    tournament_id TEXT NOT NULL,
	    matchup_id    TEXT NOT NULL,
	    FOREIGN KEY (tournament_id) REFERENCES tournaments(id),
	    FOREIGN KEY (matchup_id) REFERENCES matchups(id)
	);

	CREATE TABLE IF NOT EXISTS matchup_game (
	    id            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	    matchup_id    TEXT NOT NULL,
	    game_id       TEXT NOT NULL,
	    FOREIGN KEY (matchup_id) REFERENCES matchups(id),
	    FOREIGN KEY (game_id) REFERENCES games(id)
	);
	`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil

}

func InitDB() error {

	var err error
	db, err = sql.Open("sqlite3", "./chess.db")
	if err != nil {
		return fmt.Errorf("Failed to open database: %v", err)
	}

	err = CreateTables()
	if err != nil {
		return fmt.Errorf("Error creating tables: %v", err)
	}

	return nil

}

func CloseDB() {

	db.Close()

}

func GetDB() (*sql.DB, error) {

	if db == nil {
		return nil, errors.New("DB not initialised!")
	}
	return db, nil

}

func SafeExec(query string, args ...any) (sql.Result, error) {

	dbMutex.Lock()
	defer dbMutex.Unlock()
	return db.Exec(query, args...)

}
