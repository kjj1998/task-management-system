package models

import (
	"fmt"
	"time"
)

type DBUser struct {
	ID           uint
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u DBUser) String() string {
	return fmt.Sprintf(
		"DBUser[ID=%d, Email=%s, FirstName=%s, LastName=%s, CreatedAt=%s, UpdatedAt=%s]",
		u.ID,
		u.Email,
		u.FirstName,
		u.LastName,
		u.CreatedAt.Format(time.RFC3339),
		u.UpdatedAt.Format(time.RFC3339),
	)
}
