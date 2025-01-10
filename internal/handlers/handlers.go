package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/morzisorn/practicum_final_project/internal/logic"
)

func HandleNextDate(c *fiber.Ctx) error {
	now, err := time.Parse(logic.DateFormat, c.Query("now"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("now is invalid")
	}
	date := c.Query("date")
	repeat := c.Query("repeat")
	nextDate, err := logic.NextDate(now, date, repeat)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendString(nextDate)

}

func HandleGetTasks(c *fiber.Ctx) error {
	search := c.Query("search")
	tasks, err := logic.UpcomingTasks(search)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"tasks": tasks,
	})
}

func HandleGetTask(c *fiber.Ctx) error {
	task, err := logic.GetTask(c.Query("id"))
	if err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(task)
}

func HandleUpdateTask(c *fiber.Ctx) error {
	var req logic.Task
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

	err := logic.UpdateTask(req.ID, req.Date, req.Title, req.Comment, req.Repeat)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{})
}

func HandleAddTask(c *fiber.Ctx) error {
	var req logic.Task
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

	id, err := logic.AddTask(req.Date, req.Title, req.Comment, req.Repeat)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"id": id,
	})
}

func HandleDoneTask(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.JSON(fiber.Map{
			"error": "id is empty",
		})
	}
	err := logic.DoneTask(id)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{})
}

func HandleSignIn(c *fiber.Ctx) error {
	var req struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	token, err := logic.GetToken(req.Password)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": "password is invalid",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func HandleDeleteTask(c *fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return c.JSON(fiber.Map{
			"error": "id is empty",
		})
	}
	err := logic.DeleteTask(id)
	if err != nil {
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{})
}
