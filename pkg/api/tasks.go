package api

import (
	"encoding/json"
	"final/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler50(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // Запрашиваем максимум 50 задач
	if err != nil {
		writeJsonError(w, err.Error())
		return
	}

	// Если задач нет, создаем пустой слайс
	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func writeJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeJsonError(w http.ResponseWriter, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": errorMsg})
}
