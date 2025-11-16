package api

import (
	"encoding/json"
	"final/pkg/db"
	"final/pkg/db/nexdate"
	"fmt"
	"net/http"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		handleGetTask(w, r)
	case http.MethodPut:
		handlePutTask(w, r)
	case http.MethodDelete:
		TaskDeleteHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	if id == "" {
		WriteJson(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	// Получаем задачу по идентификатору
	task, err := db.GetTaskByID(id)
	if err != nil {
		WriteJson(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	// Возвращаем данные задачи
	WriteJson(w, task, http.StatusOK)
}

func handlePutTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var task db.Task
	if err := decoder.Decode(&task); err != nil {
		WriteJson(w, map[string]string{"error": "Ошибка десериализации JSON: " + err.Error()}, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		WriteJson(w, map[string]string{"error": "Не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	// Обработка даты с помощью функции checkDate
	if err := checkDate(&task); err != nil {
		WriteJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Обновление задачи в базе данных
	if err := db.UpdateTask(&task); err != nil {
		WriteJson(w, map[string]string{"error": "Ошибка обновления задачи: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	WriteJson(w, map[string]interface{}{"message": "Задача обновлена успешно"}, http.StatusOK)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var task db.Task
	if err := decoder.Decode(&task); err != nil {
		WriteJson(w, map[string]string{"error": "Ошибка десериализации JSON: " + err.Error()}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		WriteJson(w, map[string]string{"error": "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Обработка даты с помощью функции checkDate
	if err := checkDate(&task); err != nil {
		WriteJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Добавляем задачу в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		WriteJson(w, map[string]string{"error": "Ошибка добавления задачи: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	WriteJson(w, map[string]string{"id": fmt.Sprintf("%d", id)}, http.StatusCreated)
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJson(w, map[string]string{"error": "Не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}
	// Получаем задачу
	task, err := db.GetTaskByID(id)
	if err != nil {
		WriteJson(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}
	// Если задача не периодическая — удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			WriteJson(w, map[string]string{"error": "Ошибка удаления задачи: " + err.Error()}, http.StatusInternalServerError)
			return
		}
		WriteJson(w, map[string]interface{}{}, http.StatusOK) // Возвращаем пустой JSON-объект при успехе
		return
	}
	// Если задача периодическая — рассчитываем следующую дату
	next, err := nexdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		WriteJson(w, map[string]string{"error": "Ошибка расчёта следующей даты: " + err.Error()}, http.StatusInternalServerError)
		return
	}
	// Обновляем дату
	if err := db.UpdateDate(next, id); err != nil {
		WriteJson(w, map[string]string{"error": "Ошибка обновления даты: " + err.Error()}, http.StatusInternalServerError)
		return
	}
	WriteJson(w, map[string]interface{}{}, http.StatusOK) // Возвращаем пустой JSON-объект при успехе
}

// TaskDeleteHandler обрабатывает DELETE /api/task
func TaskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJson(w, map[string]string{"error": "Не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		WriteJson(w, map[string]string{"error": "Ошибка удаления задачи: " + err.Error()}, http.StatusInternalServerError)
		return
	}
	// Возвращаем пустой JSON-объект при успехе
	WriteJson(w, map[string]interface{}{}, http.StatusOK)
}

// Функция проверки даты
func checkDate(task *db.Task) error {
	now := time.Now() // Получаем текущее время

	// Если task.Date пустая, ставим сегодняшнюю дату
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	// Парсим дату
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("The date is presented in the wrong format")
	}

	// Проверяем, если есть правило повторения
	if task.Repeat != "" {
		next, err := nexdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("Incorrect format of the repetition rule")
		}

		// Устанавливаем следующую дату, если она больше текущей
		if t.Before(now) {
			task.Date = next
		}
	}

	// Убедимся, что дата задачи (t) не меньше текущей
	if t.Before(now) {
		task.Date = now.Format(dateFormat)
	}

	return nil
}

func WriteJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode) // Устанавливаем код состояния
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
