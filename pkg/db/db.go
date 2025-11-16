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
