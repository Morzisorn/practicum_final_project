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

	app.Get("/api/tasks", func(c *fiber.Ctx) error {
		search := c.Query("search")
		tasks, err := UpcomingTasks(search)
		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		result := make([]map[string]string, len(tasks))
		for i, task := range tasks {
			result[i] = map[string]string{
				"id":      task.ID,
				"date":    task.Date,
				"title":   task.Title,
				"comment": task.Comment,
				"repeat":  task.Repeat,
			}
		}

		return c.JSON(fiber.Map{
			"tasks": tasks,
		})
	})

	app.Get("/api/task", func(c *fiber.Ctx) error {
		task, err := GetTask(c.Query("id"))
		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(task)
	})

	app.Put("/api/task", func(c *fiber.Ctx) error {
		var req Task
		if err := c.BodyParser(&req); err != nil {
			c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if req.ID == "" {
			return c.JSON(fiber.Map{
				"error": "id is empty",
			})
		}
		if req.Comment == "" {
			req.Comment = ""
		}

		if req.Repeat == "" {
			req.Repeat = ""
		}

		err := UpdateTask(req.ID, req.Date, req.Title, req.Comment, req.Repeat)
		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{})
	})

	app.Post("/api/task", func(c *fiber.Ctx) error {
		var req Task
		if err := c.BodyParser(&req); err != nil {
			c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if req.Comment == "" {
			req.Comment = ""
		}

		if req.Repeat == "" {
			req.Repeat = ""
		}

		id, err := AddTask(req.Date, req.Title, req.Comment, req.Repeat)
		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"id": id,
		})
	})

	app.Post("/api/task/done", func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			return c.JSON(fiber.Map{
				"error": "id is empty",
			})
		}
		err := DoneTask(id)
		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{})
	})

	app.Delete("/api/task", func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			return c.JSON(fiber.Map{
				"error": "id is empty",
			})
		}
		err := DeleteTask(id)
		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{})
	})

	port := os.Getenv("TODO_PORT")

	logrus.Fatal(app.Listen(fmt.Sprintf(":%s", port)))

}
