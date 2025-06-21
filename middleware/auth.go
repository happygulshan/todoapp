package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"
	"todoapp/models"
)

type key string

const userIDKey key = "userID"

func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := strings.TrimSpace(r.Header.Get("Authorization"))
			if token == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			var session models.Session
			err := db.QueryRow("SELECT * FROM sessions WHERE token = $1", token).
				Scan(&session.ID, &session.User_id, &session.Token, &session.CreatedAt, &session.Expires_at)

			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if session.Expires_at.Before(time.Now()) {
				http.Error(w, "Session expired", http.StatusUnauthorized)
				return
			}

			// Add userID to request context
			ctx := context.WithValue(r.Context(), userIDKey, session.User_id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper to get userID from context
func GetUserID(r *http.Request) string {
	id := r.Context().Value(userIDKey)
	if idStr, ok := id.(string); ok {
		return idStr
	}
	return ""
}
