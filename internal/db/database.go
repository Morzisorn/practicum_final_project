package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/morzisorn/practicum_final_project/config"
	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

var (
	DB *sql.DB
)

func createDBTables() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date VARCHAR(8),
		title VARCHAR(255) NOT NULL,
		comment TEXT,
		repeat VARCHAR(255)
	);`)
	if err != nil {
		logrus.Fatal(err)
	}

	_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);`)
	if err != nil {
		logrus.Fatal(err)
	}
}

func RunDB() error {
	appPath, err := os.Getwd()
	if err != nil {
		return err
	}

	DBFile := filepath.Join(appPath, config.DBFilename)
	_, err = os.Stat(DBFile)

	var install bool
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", DBFile)
	if err != nil {
		return err
	}

	if install {
		createDBTables()
	}
	return nil
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			logrus.Fatalf("DB close error: %v", err)
		}
	}
}
