package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/morzisorn/practicum_final_project/internal/db"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const (
	DateFormat         = "20060102"
	upcomingTasksLimit = 10
)

func CheckDateAndRepeat(date string, repeat string) (string, error) {
	if date == "" {
		date = time.Now().Format(DateFormat)
	} else {
		dateTime, err := time.Parse(DateFormat, date)
		if err != nil {
			return "", fmt.Errorf("date is invalid")
		}
		if dateTime.Before(time.Now()) && repeat == "" {
			date = time.Now().Format(DateFormat)
		}
	}

	return date, nil
}

func AddTask(date, title, comment, repeat string) (int64, error) {
	date, err := CheckDateAndRepeat(date, repeat)
	if err != nil {
		return 0, err
	}

	if repeat != "" && date < time.Now().Format(DateFormat) {
		date, err = NextDate(time.Now(), date, repeat)
		if err != nil {
			return 0, fmt.Errorf("repeat is invalid")
		}
	}

	if title == "" {
		return 0, fmt.Errorf("title is empty")
	}

	var id int
	if db.DB == nil {
		return 0, fmt.Errorf("db is nil")
	}
	err = db.DB.QueryRowContext(context.Background(), "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?) RETURNING id", date, title, comment, repeat).Scan(&id)
	if err != nil {
		return 0, err
	}
	return int64(id), nil

}

func UpcomingTasks(search string) ([]Task, error) {
	if search == "" {
		query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER by date LIMIT ?"
		return getTasksFromDB(query, upcomingTasksLimit)
	}

	date, err := time.Parse("02.01.2006", search)
	//Search by date
	if err == nil {
		query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER by date LIMIT ?"
		return getTasksFromDB(query, date.Format(DateFormat), upcomingTasksLimit)
	}

	//Search by title and comment
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"
	return getTasksFromDB(query, "%"+search+"%", "%"+search+"%", upcomingTasksLimit)
}

func getTasksFromDB(query string, args ...interface{}) ([]Task, error) {
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return []Task{}, err
	}

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []Task{}, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return []Task{}, err
	}
	if err = rows.Close(); err != nil {
		return []Task{}, err
	}
	return tasks, nil
}

func GetTask(id string) (Task, error) {
	var task Task
	err := db.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func UpdateTask(id, date, title, comment, repeat string) error {
	date, err := CheckDateAndRepeat(date, repeat)
	if err != nil {
		return err
	}

	if repeat != "" {
		var err error
		date, err = NextDate(time.Now(), date, repeat)
		if err != nil {
			return fmt.Errorf("repeat is invalid")
		}
	}

	if title == "" {
		return fmt.Errorf("title is empty")
	}

	res, err := db.DB.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", date, title, comment, repeat, id)
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

func DoneTask(id string) error {
	task, err := GetTask(id)
	if err != nil {
		return err
	}
	if task.Repeat == "" {
		_, err := db.DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			return err
		}
		return nil
	}

	err = UpdateTask(id, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := db.DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
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
