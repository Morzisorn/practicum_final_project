package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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

	appPath, err := os.Executable()
	if err != nil {
		logrus.Fatal(err)
	}
	filename := os.Getenv("TODO_DBFILE")
	dbFile := filepath.Join(filepath.Dir(appPath), filename)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		logrus.Fatal(err)
	}

	if install {
		createDBTables()
	} else {
		// Проверяем, существует ли таблица
		row := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='scheduler';")
		var name string
		if err := row.Scan(&name); err != nil {
			fmt.Println("Table 'scheduler' does not exist in the database. Creating again.")
			createDBTables()
		}
		fmt.Println("Table 'scheduler' already exists.")
	}

	startServer()
}
