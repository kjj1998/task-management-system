package models

import "time"

type DbUser struct {
	ID           uint
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
