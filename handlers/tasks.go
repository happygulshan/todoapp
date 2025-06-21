package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todoapp/middleware"
	"todoapp/models"

	"github.com/gorilla/mux"

	// "todoapp/models"

	// "todoapp/middleware"
	"strings"
)

// handlers/task.go
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {

	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r)

	task.User_id = userID
	if task.Status == "" {
		task.Status = "pending"
	}

	err := h.DB.QueryRow("INSERT INTO tasks(title, description, user_id, status) VALUES($1, $2, $3, $4) RETURNING id",
		task.Title, task.Description, task.User_id, task.Status).Scan(&task.ID)

	if err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]
	var user_id string
	err := h.DB.QueryRow("select user_id from tasks where id = $1", taskID).Scan(&user_id)

	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	sessionUserId := middleware.GetUserID(r)
	if user_id != sessionUserId {
		http.Error(w, "you dont have access to update it", http.StatusForbidden)
		return
	}

	var updateFields map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateFields); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	allowed := map[string]bool{"title": true, "description": true, "status": true}

	for key := range updateFields {
		if !allowed[key] {
			http.Error(w, "invalid field: "+key, http.StatusBadRequest)
			return
		}
	}

	query := "UPDATE tasks SET "
	params := []interface{}{}
	i := 1

	for key, value := range updateFields {
		query += fmt.Sprintf("%s = $%d, ", key, i)
		params = append(params, value)
		i++
	}

	// Remove last comma
	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = $%d" // for task ID
	query = fmt.Sprintf(query, i)

	params = append(params, taskID)

	var task models.Task
	query += " RETURNING id, title, description, user_id, status"

	err = h.DB.QueryRow(query, params...).Scan(&task.ID, &task.Title, &task.Description, &task.User_id, &task.Status)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "failed to update task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)
	rows, err := h.DB.Query("SELECT id, title, description, user_id, status FROM tasks WHERE user_id = $1::uuid", userID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.User_id, &task.Status); err != nil {
			http.Error(w, "failed to scan task", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	// Handle case where no rows matched
	if len(tasks) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "No tasks found"})
		return
	}

	// Send back the tasks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	taskID := vars["id"]
	var user_id string
	err := h.DB.QueryRow("select user_id from tasks where id = $1", taskID).Scan(&user_id)

	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	sessionUserId := middleware.GetUserID(r)

	if user_id != sessionUserId {
		http.Error(w, "you dont have access to delete it", http.StatusUnauthorized)
		return
	}

	_, err = h.DB.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"msg": "Successfully deleted task",
	})
}
