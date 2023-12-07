package db

import "time"

type Score struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}
