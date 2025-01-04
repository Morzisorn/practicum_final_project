package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	app.Get("/api/nextdate", func(c *fiber.Ctx) error {
		now, err := time.Parse("20060102", c.Query("now"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("now is invalid")
		}
		date := c.Query("date")
		repeat := c.Query("repeat")
		nextDate, err := NextDate(now, date, repeat)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendString(nextDate)

	})

	port := os.Getenv("TODO_PORT")

	logrus.Fatal(app.Listen(fmt.Sprintf(":%s", port)))

}
