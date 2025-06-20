package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	//basic simple validation
	if user.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if user.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	// Validate email format (simple regex)
	if !isValidEmail(user.Email) {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}

	// Validate password strength
	if len(user.Password) < 8 {
		http.Error(w, "password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	var existingID string
	fmt.Println(user.Email)
	err := h.DB.QueryRow("SELECT id FROM users WHERE email = TRIM($1)", user.Email).Scan(&existingID)

	if err == nil {
		http.Error(w, "email already exists", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	// end of validation

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Save user to DB
	err = h.DB.QueryRow("INSERT INTO users(name, email, password) VALUES($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password).Scan(&user.ID)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "faiiled to register user", http.StatusInternalServerError)
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
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	var hashedPassword string

	// Only fetch hashed password
	err := h.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email=$1", req.Email).
		Scan(&user.ID, &user.Name, &user.Email, &hashedPassword)

	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the hashed one
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate session token
	token := uuid.New().String()

	// Insert session into DB
	_, err = h.DB.Exec("INSERT INTO sessions(user_id, token) VALUES($1, $2)", user.ID, token)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	// Return token to client
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"msg":   "Login Successful",
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.Header.Get("Authorization"))

	var expires_at time.Time
	// check if token already expired

	err := h.DB.QueryRow("SELECT expires_at FROM sessions WHERE token = $1", token).
		Scan(&expires_at)

	if err != nil {
		http.Error(w, "something wrong in veryfying token", http.StatusUnauthorized)
		return
	}

	if expires_at.Before(time.Now()) {
		http.Error(w, "session expired. Please log in again.", http.StatusUnauthorized)
		return
	}

	var tok string
	err = h.DB.QueryRow("UPDATE sessions SET expires_at = CURRENT_TIMESTAMP WHERE token=$1 RETURNING token", token).Scan(&tok)

	if err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}

	// Return token and msg to client
	json.NewEncoder(w).Encode(map[string]string{
		"token": tok,
		"msg":   "Successful logout",
	})

}
