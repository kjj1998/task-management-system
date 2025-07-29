package models

import (
	"fmt"
	"time"
)

type (
	TaskStatus   string
	TaskPriority string
)

const (
	Low    TaskPriority = "low"
	Medium TaskPriority = "medium"
	High   TaskPriority = "high"
)

const (
	Pending    TaskStatus = "pending"
	InProgress TaskStatus = "in_progress"
	Completed  TaskStatus = "completed"
)

type DBTask struct {
	ID          string       `json:"id"`
	UserID      string       `json:"userID"`
	CategoryID  string       `json:"categoryID"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	Status      TaskStatus   `json:"status"`
	DueDate     *time.Time   `json:"dueDate"`
	CompletedAt *time.Time   `json:"completedAt"`
	CreatedAt   *time.Time   `json:"createdAt"`
	UpdatedAt   *time.Time   `json:"updatedAt"`
}

func (t DBTask) String() string {
	formatTime := func(t *time.Time) string {
		if t == nil {
			return "nil"
		}
		return t.Format(time.RFC3339)
	}

	return fmt.Sprintf(
		"DBTask[ID=%s, UserID=%s, CategoryID=%s, Title=%s, Description=%s, Priority=%s, Status=%s, DueDate=%s, CompletedAt=%s, CreatedAt=%s, UpdatedAt=%s]",
		t.ID,
		t.UserID,
		t.CategoryID,
		t.Title,
		t.Description,
		t.Priority,
		t.Status,
		formatTime(t.DueDate),
		formatTime(t.CompletedAt),
		formatTime(t.CreatedAt),
		formatTime(t.UpdatedAt),
	)
}
