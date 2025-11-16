package api

import (
	"encoding/json"
	"final/pkg/db"
	"net/http"
)

const TaskLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(TaskLimit) // Запрашиваем максимум 50 задач
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
	}, http.StatusOK) // Добавлена запятая и исправлен статус-код
}

func writeJson(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode) // Устанавливаем код состояния
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJsonError(w http.ResponseWriter, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError) // Добавлен статус-код для ошибки
	json.NewEncoder(w).Encode(map[string]string{"error": errorMsg})
}
