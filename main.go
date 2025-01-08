package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func startServer() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 600 * time.Second,
	})

	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./web/index.html")
	})

	app.Get("/index.html", func(c *fiber.Ctx) error {
		return c.SendFile("./web/index.html")
	})

	app.Get("/login.html", func(c *fiber.Ctx) error {
		return c.SendFile("./web/login.html")
	})

	app.Get("/js/axios.min.js", func(c *fiber.Ctx) error {
		return c.SendFile("./web/js/axios.min.js")
	})

	app.Get("/js/scripts.min.js", func(c *fiber.Ctx) error {
		return c.SendFile("./web/js/scripts.min.js")
	})

	app.Get("/css/style.css", func(c *fiber.Ctx) error {
		return c.SendFile("./web/css/style.css")
	})

	app.Get("/css/theme.css", func(c *fiber.Ctx) error {
		return c.SendFile("./web/css/theme.css")
	})

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendFile("./web/favicon.ico")
	})

	app.Get("/api/nextdate", handleNextDate)

	app.Get("/api/tasks", auth, handleGetTasks)

	app.Get("/api/task", auth, handleGetTask)

	app.Put("/api/task", auth, handleUpdateTask)

	app.Post("/api/task", auth, handleAddTask)

	app.Post("/api/task/done", auth, handleDoneTask)

	app.Post("/api/signin", handleSignIn)

	app.Delete("/api/task", auth, handleDeleteTask)

	port := os.Getenv("TODO_PORT")

	logrus.Fatal(app.Listen(fmt.Sprintf(":%s", port)))

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
