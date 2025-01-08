package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func handleNextDate(c *fiber.Ctx) error {
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

}

func handleGetTasks(c *fiber.Ctx) error {
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
}

func handleGetTask(c *fiber.Ctx) error {
	task, err := GetTask(c.Query("id"))
	if err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(task)
}

func handleUpdateTask(c *fiber.Ctx) error {
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
}

func handleAddTask(c *fiber.Ctx) error {
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
}

func handleDoneTask(c *fiber.Ctx) error {
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
}

func handleSignIn(c *fiber.Ctx) error {
	var req struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	token, err := getToken(req.Password)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": "password is invalid",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func handleDeleteTask(c *fiber.Ctx) error {
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
}
