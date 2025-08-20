package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(255) NOT NULL DEFAULT ""
);

CREATE INDEX idx_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	if install {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to initialize database schema: %w", err)
		}
	}

	return nil
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `
	INSERT INTO scheduler (date, title, comment, repeat) VALUES
		(:date, :title, :comment, :repeat);
	`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}
