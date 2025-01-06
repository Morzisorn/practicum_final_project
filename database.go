package main

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

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
