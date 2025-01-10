package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/mattn/go-sqlite3"
	"github.com/morzisorn/practicum_final_project/config"
	"github.com/morzisorn/practicum_final_project/internal/db"
	"github.com/morzisorn/practicum_final_project/internal/handlers"
	"github.com/morzisorn/practicum_final_project/internal/logic"
	"github.com/sirupsen/logrus"
)

func startServer() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
	})

	app.Use(recover.New())

	app.Static("/", "./web")

	app.Get("/api/nextdate", handlers.HandleNextDate)

	app.Get("/api/tasks", logic.Auth, handlers.HandleGetTasks)

	app.Get("/api/task", logic.Auth, handlers.HandleGetTask)

	app.Put("/api/task", logic.Auth, handlers.HandleUpdateTask)

	app.Post("/api/task", logic.Auth, handlers.HandleAddTask)

	app.Post("/api/task/done", logic.Auth, handlers.HandleDoneTask)

	app.Post("/api/signin", handlers.HandleSignIn)

	app.Delete("/api/task", logic.Auth, handlers.HandleDeleteTask)

	logrus.Fatal(app.Listen(fmt.Sprintf(":%s", config.Port)))

}

func main() {
	err := config.GetEnvs()
	if err != nil {
		logrus.Fatalf("getEnvs error: %v", err)
	}

	defer db.CloseDB()
	err = db.RunDB()
	if err != nil {
		logrus.Fatalf("RunDB error: %v", err)
	}

	startServer()
}
