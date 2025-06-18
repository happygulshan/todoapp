package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"todoapp/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB *sql.DB
}

// Simple email validation but need complex regex for production
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	//basic simple validation
	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if user.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Validate email format (simple regex)
	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate password strength
	if len(user.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	var existingID string
	fmt.Println(user.Email)
	err := h.DB.QueryRow("SELECT id FROM users WHERE email = TRIM($1)", user.Email).Scan(&existingID)

	if err == nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	// end of validation

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Save user to DB
	err = h.DB.QueryRow("INSERT INTO users(name, email, password) VALUES($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password).Scan(&user.ID)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Faiiled to register user", http.StatusInternalServerError)
		return
	}

	//password not sending back
	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	var hashedPassword string

	// Only fetch hashed password
	err := h.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email=$1", req.Email).
		Scan(&user.ID, &user.Name, &user.Email, &hashedPassword)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the hashed one
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate session token
	token := uuid.New().String()

	// Insert session into DB
	_, err = h.DB.Exec("INSERT INTO sessions(user_id, token) VALUES($1, $2)", user.ID, token)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Return token to client
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.Header.Get("Authorization"))
	fmt.Println(token)
}
