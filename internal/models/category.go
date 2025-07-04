package models

import (
	"fmt"
	"time"
)

type DBCategory struct {
	ID        string
	UserID    string
	Name      string
	Color     string
	CreatedAt *time.Time
}

func (c DBCategory) String() string {
	return fmt.Sprintf(
		"DBCategory[ID=%s, UserID=%s, Name=%s, Color=%s, CreatedAt=%s]",
		c.ID,
		c.UserID,
		c.Name,
		c.Color,
		c.CreatedAt.Format(time.RFC3339),
	)
}
