package api

import (
	"net/http"

	"github.com/Solbreen/final-project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, TaskResp{Error: "ID parameter is required"}, http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, TaskResp{Error: err.Error()}, http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]string{}, http.StatusOK)
}
