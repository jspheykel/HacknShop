package models

import "time"

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	IsAdmin      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
