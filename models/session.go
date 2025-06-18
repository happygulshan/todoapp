package models

import "time"

type Session struct {
	ID         string    `json:"id"`
	User_id    string    `json:"user_id"`
	Token      string    `json:"token"`
	CreatedAt  time.Time `json:"created_at"`
	Expires_at time.Time `json:"expires_at"`
}
