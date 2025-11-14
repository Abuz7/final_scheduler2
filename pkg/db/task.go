package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	return id, err
}
func isValidDate(dateStr string) bool {
	if len(dateStr) != 8 {
		return false
	}
	_, err := time.Parse("20060102", dateStr)
	return err == nil
}

func isValidTitle(title string) bool {
	return len(strings.TrimSpace(title)) > 0
}

func isValidRepeat(repeat string) bool {
	if repeat == "" {
		return true
	}
	parts := strings.Fields(repeat)
	if len(parts) != 2 {
		return false
	}
	_, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	return parts[0] == "d" || parts[0] == "w" || parts[0] == "m"
}

func UpdateTask(task *Task) error {
	// Проверяем корректность даты
	if !isValidDate(task.Date) {
		return fmt.Errorf("некорректный формат даты: %s", task.Date)
	}
	// Проверяем корректность заголовка
	if !isValidTitle(task.Title) {
		return fmt.Errorf("заголовок не может быть пустым")
	}
	// Проверяем корректность повтора
	if !isValidRepeat(task.Repeat) {
		return fmt.Errorf("некорректный формат повтора: %s", task.Repeat)
	}
	// Проверяем, существует ли задача с указанным идентификатором
	var count int
	queryCheck := `SELECT COUNT(*) FROM scheduler WHERE id = ?`
	if err := GetDB().QueryRow(queryCheck, task.ID).Scan(&count); err != nil {
		return err
	}
	// Если задача не найдена, возвращаем ошибку
	if count == 0 {
		return fmt.Errorf("некорректный идентификатор для обновления задачи")
	}
	// Проверяем, есть ли изменения в данных
	queryCheckNoChange := `SELECT COUNT(*) FROM scheduler
	                        WHERE id = ? AND date = ? AND title = ? AND comment = ? AND repeat = ?`
	var noChangeCount int
	err := GetDB().QueryRow(queryCheckNoChange, task.ID, task.Date, task.Title, task.Comment, task.Repeat).Scan(&noChangeCount)
	if err != nil {
		return err
	}
	// Если данные совпадают, возвращаем ошибку
	if noChangeCount > 0 {
		return fmt.Errorf("нет изменений для задачи с идентификатором: %s", task.ID)
	}
	// Обновляем задачу
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	// Проверяем количество затронутых строк
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return fmt.Errorf("задача с идентификатором %s была найдена, но не обновлена", task.ID)
	}
	return nil
}

func GetTaskByID(id string) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := GetDB().QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
