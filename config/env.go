package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Port       string
	DBFilename string
	RealPass   string
)

func GetEnvs() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	envPath := fmt.Sprintf("%s/config/.env", dir)
	err = godotenv.Load(envPath)
	if err != nil {
		fmt.Println(envPath)
		fmt.Println("Load .env error")
	}
	Port = os.Getenv("TODO_PORT")
	DBFilename = os.Getenv("TODO_DBFILE")
	RealPass = os.Getenv("TODO_PASSWORD")
	return nil
}
