package main

import (
	"database/sql"
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Load .env error")
	}

	defer CloseDB()
	err = RunDB()
	if err != nil {
		logrus.Fatalf("RunDB error: %v", err)
	}

	startServer()
}
