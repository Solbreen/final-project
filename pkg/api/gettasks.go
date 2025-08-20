package api

import (
	"net/http"
	"time"

	"github.com/Solbreen/final-project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
	Error string     `json:"error,omitempty"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {

	var (
		err   error
		tasks []*db.Task
		limit int = 50
	)

	search := r.URL.Query().Get("search")

	if search != "" {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			dateStr := date.Format("20060102")
			tasks, err = db.GetTasksByDate(dateStr, limit)
		} else {
			tasks, err = db.GetTasksBySearch(search, limit)
		}
	} else {
		tasks, err = db.Tasks(limit)
	}

	if err != nil {
		writeJSON(w, TasksResp{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}
	writeJSON(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
