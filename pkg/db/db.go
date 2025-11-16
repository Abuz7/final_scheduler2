package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL,
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);
CREATE INDEX idx_date ON scheduler(date);
`

// Init initializes the SQLite database
func Init(dbFile string) (*sql.DB, error) {
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		// Файл не существует, создадим базу данных
		db, err = sql.Open("sqlite", dbFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %v", err)
		}
		_, err = db.Exec(schema)
		if err != nil {
			return nil, fmt.Errorf("failed to execute schema: %v", err)
		}
	} else {
		// Файл существует, просто открываем базу
		db, err = sql.Open("sqlite", dbFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %v", err)
		}
	}
	return db, nil
}

// GetDB returns the database reference
func GetDB() *sql.DB {
	return db
}
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
	rows, err := GetDB().Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := GetDB().Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %s not found", id)
	}

	return nil
}

// UpdateDate обновляет только дату задачи
func UpdateDate(nextDate string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err := GetDB().Exec(query, nextDate, id)
	return err
}
