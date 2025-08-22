package api

import (
	"net/http"

	"github.com/Solbreen/final-project/pkg/db"
)

type TaskResp struct {
	Task  *db.Task `json:"task"`
	Error string   `json:"error,omitempty"`
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, TaskResp{Error: "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, TaskResp{Error: "ID parameter is required"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, TaskResp{Error: err.Error()}, http.StatusNotFound)
		return
	}
	writeJSON(w, task, http.StatusOK)
}
