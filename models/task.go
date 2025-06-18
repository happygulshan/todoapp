package models

import "time"

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	User_id     string    `json:"user_id"`
	Status      string    `json:"status"`
	Created_at  time.Time `json:"created_at"`
}
