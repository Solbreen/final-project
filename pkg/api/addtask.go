package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Solbreen/final-project/pkg/db"
)

type TaskResponse struct {
	ID    int64  `json:"id"`
	Error string `json:"error,omitempty"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, TaskResponse{Error: "Ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSON(w, TaskResponse{Error: "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, TaskResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, TaskResponse{ID: id}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}
	if t.Before(now) {
		if task.Repeat != "" {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}

			if afterNow(now, t) {
				task.Date = next
			}
		} else {
			task.Date = now.Format("20060102")
		}
	}
	return nil
}
