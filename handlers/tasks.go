package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todoapp/models"

	// "todoapp/models"

	// "todoapp/middleware"
	"strings"
)

// handlers/task.go
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {

	// userID := middleware.GetUserID(r)
	token := strings.TrimSpace(r.Header.Get("Authorization"))

	var session models.Session

	// var tok string
	err := h.DB.QueryRow("SELECT * FROM sessions WHERE token = $1", token).
		Scan(&session.ID, &session.User_id, &session.Token, &session.CreatedAt, &session.Expires_at)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if session.Expires_at.Before(time.Now()) {
		http.Error(w, "Session expired. Please log in again.", http.StatusUnauthorized)
		return
	}

	var task models.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	task.User_id = session.User_id
	if task.Status == "" {
		task.Status = "pending"
	}

	err = h.DB.QueryRow("INSERT INTO tasks(title, description, user_id, status) VALUES($1, $2, $3, $4) RETURNING id",
		task.Title, task.Description, task.User_id, task.Status).Scan(&task.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}
