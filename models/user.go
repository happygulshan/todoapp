package models

import "time"

type User struct {
	ID        string       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`   // Make sure to hash this before storing
	CreatedAt time.Time `json:"created_at"` // Maps to PostgreSQL TIMESTAMP
}
