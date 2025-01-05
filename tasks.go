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

func AddTask(date, title, comment, repeat string) (int64, error) {
	if date == "" {
		date = time.Now().Format("20060102")
	} else {
		dateTime, err := time.Parse("20060102", date)
		if err != nil {
			return 0, fmt.Errorf("date is invalid")
		}
		if dateTime.Before(time.Now()) && repeat == "" {
			date = time.Now().Format("20060102")
		}
	}
	if repeat != "" {
		var err error
		date, err = NextDate(time.Now(), date, repeat)
		if err != nil {
			return 0, fmt.Errorf("repeat is invalid")
		}
	}

	if title == "" {
		return 0, fmt.Errorf("title is empty")
	}
	var id int
	//db := GetDB()
	if db == nil {
		return 0, fmt.Errorf("db is nil")
	}
	err := db.QueryRowContext(context.Background(), "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?) RETURNING id", date, title, comment, repeat).Scan(&id)
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
