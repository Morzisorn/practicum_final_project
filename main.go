package main

import (
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Load .env error")
	}
	createDB()
	startServer()
	defer db.Close()
}
