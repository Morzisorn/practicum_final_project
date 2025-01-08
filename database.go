package main

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	db *sql.DB
)

func createDBTables() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date VARCHAR(8),
		title VARCHAR(255) NOT NULL,
		comment TEXT,
		repeat VARCHAR(255)
	);`)
	if err != nil {
		logrus.Fatal(err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);`)
	if err != nil {
		logrus.Fatal(err)
	}
}

func RunDB() error {
	appPath, err := os.Getwd()
	if err != nil {
		return err
	}
	filename := os.Getenv("TODO_DBFILE")
	dbFile := filepath.Join(appPath, filename)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	if install {
		createDBTables()
	}
	return nil
}

func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			logrus.Fatalf("db close error: %v", err)
		}
	}
}
