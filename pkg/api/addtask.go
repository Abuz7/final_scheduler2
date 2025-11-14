package api

import (
	"encoding/json"
	"final/pkg/api/nexdate"
	"final/pkg/db"
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
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// Получаем задачу по идентификатору
	task, err := db.GetTaskByID(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	// Возвращаем данные задачи
	writeJson(w, task)
}

func handlePutTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var task db.Task
	if err := decoder.Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "Ошибка десериализации JSON: " + err.Error()})
		return
	}

	if task.ID == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	// Обработка даты с помощью функции checkDate
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Обновление задачи в базе данных
	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": "Ошибка обновления задачи: " + err.Error()})
		return
	}

	writeJson(w, map[string]interface{}{"message": "Задача обновлена успешно"})
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var task db.Task
	if err := decoder.Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "Ошибка десериализации JSON: " + err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Обработка даты с помощью функции checkDate
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Добавляем задачу в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка добавления задачи: " + err.Error()})
		return
	}

	writeJson(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}
	// Получаем задачу
	task, err := db.GetTaskByID(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}
	// Если задача не периодическая — удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": "Ошибка удаления задачи: " + err.Error()})
			return
		}
		writeJson(w, map[string]interface{}{})
		return
	}
	// Если задача периодическая — рассчитываем следующую дату
	next, err := nexdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка расчёта следующей даты: " + err.Error()})
		return
	}
	// Обновляем дату
	if err := db.UpdateDate(next, id); err != nil {
		writeJson(w, map[string]string{"error": "Ошибка обновления даты: " + err.Error()})
		return
	}
	writeJson(w, map[string]interface{}{})
}

// TaskDeleteHandler обрабатывает DELETE /api/task
func TaskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": "Ошибка удаления задачи: " + err.Error()})
		return
	}
	// Возвращаем пустой JSON-объект при успехе
	writeJson(w, map[string]interface{}{})
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
		return fmt.Errorf("Дата представлена в неверном формате")
	}

	// Проверяем, если есть правило повторения
	if task.Repeat != "" {
		next, err := nexdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("Неправильный формат правила повторения")
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

func WriteJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
