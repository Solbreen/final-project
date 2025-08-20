package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat"`
}

func Tasks(limit int) ([]*Task, error) {
	query := `
        SELECT id, date, title, comment, repeat 
        FROM scheduler 
        WHERE date >= date('now')
        ORDER BY date ASC
        LIMIT :limit
    `

	rows, err := db.Query(query, sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tasks, nil
}

func GetTasksBySearch(search string, limit int) ([]*Task, error) {
	searchPattern := "%" + search + "%"
	query := `
        SELECT id, date, title, comment, repeat 
        FROM scheduler 
        WHERE (lower(title) LIKE lower(:sp1) OR lower(comment) LIKE lower(:sp2))
		AND date >= strftime('%Y%m%d', 'now')
        ORDER BY date ASC
        LIMIT :limit
    `

	rows, err := db.Query(query, sql.Named("sp1", searchPattern), sql.Named("sp2", searchPattern), sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func GetTasksByDate(date string, limit int) ([]*Task, error) {
	query := `
        SELECT id, date, title, comment, repeat 
        FROM scheduler 
        WHERE date = :date
        ORDER BY date ASC
        LIMIT :limit
    `

	rows, err := db.Query(query, sql.Named("date", date), sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	query := `
        SELECT id, date, title, comment, repeat 
        FROM scheduler 
        WHERE id = :id
    `
	var task Task
	err := db.QueryRow(query, sql.Named("id", id)).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler 
		SET date = :date, title = :title, comment = :comment, repeat = :repeat
        WHERE id = :id
		`
	res, err := db.Exec(
		query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func UpdateDate(next, id string) error {
	query := `
		UPDATE scheduler 
		SET date = :date
        WHERE id = :id
		`
	res, err := db.Exec(
		query,
		sql.Named("date", next),
		sql.Named("id", id),
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	query := `
        DELETE FROM scheduler 
        WHERE id = :id
    `
	res, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf(`can't delete task`)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
