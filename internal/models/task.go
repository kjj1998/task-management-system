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

// func (p TaskPriority) String() string {
// 	switch p {
// 	case Low:
// 		return "low"
// 	case Medium:
// 		return "medium"
// 	case High:
// 		return "high"
// 	default:
// 		return "medium"
// 	}
// }

// func (s TaskStatus) String() string {
// 	switch s {
// 	case Pending:
// 		return "pending"
// 	case InProgress:
// 		return "in_progress"
// 	case Completed:
// 		return "completed"
// 	default:
// 		return "pending"
// 	}
// }

// func (s TaskStatus) Value() (driver.Value, error) {
// 	return s.String(), nil
// }

// func ParseTaskStatus(s string) (TaskStatus, error) {
// 	switch s {
// 	case "pending":
// 		return Pending, nil
// 	case "in_progress":
// 		return InProgress, nil
// 	case "completed":
// 		return Completed, nil
// 	default:
// 		return Pending, fmt.Errorf("invalid task status: %s", s)
// 	}
// }

// func (s *TaskStatus) Scan(value any) error {
// 	var str string

// 	switch v := value.(type) {
// 	case string:
// 		str = v
// 	case []byte:
// 		str = string(v)
// 	default:
// 		return fmt.Errorf("invalid type for TaskStatus: %T", value)
// 	}

// 	parsed, err := ParseTaskStatus(str)
// 	if err != nil {
// 		return err
// 	}

// 	*s = parsed
// 	return nil
// }

// func (s TaskPriority) Value() (driver.Value, error) {
// 	return s.String(), nil
// }

// func ParsePriorityStatus(s string) (TaskPriority, error) {
// 	switch s {
// 	case "low":
// 		return Low, nil
// 	case "medium":
// 		return Medium, nil
// 	case "high":
// 		return High, nil
// 	default:
// 		return Medium, fmt.Errorf("invalid priority status: %s", s)
// 	}
// }

// func (s *TaskPriority) Scan(value any) error {
// 	var str string

// 	switch v := value.(type) {
// 	case string:
// 		str = v
// 	case []byte:
// 		str = string(v)
// 	default:
// 		return fmt.Errorf("invalid type for TaskPriority: %T", value)
// 	}

// 	parsed, err := ParsePriorityStatus(str)
// 	if err != nil {
// 		return err
// 	}

// 	*s = parsed
// 	return nil
// }

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

type CreateTaskResponse struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}
