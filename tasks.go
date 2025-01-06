package main

import (
	"context"
	"fmt"
	"time"
)

type (
	Task struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}
)

func CheckDateAndRepeat(date *string, repeat string) error {
	if *date == "" {
		*date = time.Now().Format("20060102")
	} else {
		dateTime, err := time.Parse("20060102", *date)
		if err != nil {
			return fmt.Errorf("date is invalid")
		}
		if dateTime.Before(time.Now()) && repeat == "" {
			*date = time.Now().Format("20060102")
		}
	}
	if repeat != "" {
		var err error
		*date, err = NextDate(time.Now(), *date, repeat)
		if err != nil {
			return fmt.Errorf("repeat is invalid")
		}
	}
	return nil
}

func AddTask(date, title, comment, repeat string) (int64, error) {
	err := CheckDateAndRepeat(&date, repeat)
	if err != nil {
		return 0, err
	}

	if title == "" {
		return 0, fmt.Errorf("title is empty")
	}

	var id int
	if db == nil {
		return 0, fmt.Errorf("db is nil")
	}
	err = db.QueryRowContext(context.Background(), "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?) RETURNING id", date, title, comment, repeat).Scan(&id)
	if err != nil {
		return 0, err
	}
	return int64(id), nil

}

func UpcomingTasks(search string) ([]Task, error) {
	if search == "" {
		query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER by date LIMIT 10"
		return getTasksFromDB(query, search)
	}

	date, err := time.Parse("02.01.2006", search)
	//Search by date
	if err == nil {
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER by date LIMIT 10"
		return getTasksFromDB(query, date.Format("20060102"))
	}

	//Search by title and comment
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 10"
	return getTasksFromDB(query, "%"+search+"%", "%"+search+"%")
}

func getTasksFromDB(query string, args ...interface{}) ([]Task, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()
	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []Task{}, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTask(id string) (Task, error) {
	var task Task
	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func UpdateTask(id, date, title, comment, repeat string) error {
	err := CheckDateAndRepeat(&date, repeat)
	if err != nil {
		return err
	}

	if title == "" {
		return fmt.Errorf("title is empty")
	}

	res, err := db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", date, title, comment, repeat, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
