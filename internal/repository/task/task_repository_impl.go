package task

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kjj1998/task-management-system/internal/models"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (t *taskRepository) Create(task *models.DBTask) (string, error) {
	command := "INSERT INTO tasks (id, user_id, category_id, title, description, priority, status, due_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	task_id := uuid.NewString()
	_, err := t.db.Exec(
		command,
		task_id,
		task.UserID,
		task.CategoryID,
		task.Title,
		task.Description,
		task.Priority.String(),
		task.Status.String(),
		task.DueDate,
	)
	if err != nil {
		return "", fmt.Errorf("createTask: %v", err)
	}

	return task_id, nil
}

func (t *taskRepository) GetAllForUser(user_id string) ([]models.DBTask, error) {
	query := "SELECT id, user_id, category_id, title, description, priority, status, due_date, completed_at, created_at, updated_at FROM tasks WHERE user_id = ?"
	rows, err := t.db.Query(query, user_id)
	if err != nil {
		return nil, fmt.Errorf("tasksByUser %q: %v", user_id, err)
	}

	var tasks []models.DBTask
	for rows.Next() {
		var task models.DBTask
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.CategoryID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.DueDate,
			&task.CompletedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("tasksByUser %q: %v", user_id, err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("tasksByUser %q: %v", user_id, err)
	}

	return tasks, nil
}

func (t *taskRepository) GetById(task_id string) (task models.DBTask, err error) {
	query := "SELECT * FROM tasks WHERE id = ?"
	row := t.db.QueryRow(query, task_id)
	err = row.Scan(
		&task.ID,
		&task.UserID,
		&task.CategoryID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Status,
		&task.DueDate,
		&task.CompletedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, fmt.Errorf("taskById %s: no such task", task_id)
		}
		return task, fmt.Errorf("taskById %s: %v", task_id, err)
	}

	return task, nil
}

func (t *taskRepository) Update(task *models.DBTask) error {
	// TODO: Write a function that forms the SQL query when given an input of string[string] map
	command := "UPDATE tasks SET " +
		"title = ?, " +
		"description = ?, " +
		"priority = ?, " +
		"status = ?, " +
		"due_date = ?, " +
		"completed_at = ?, " +
		"updated_at = ? " +
		"WHERE id = ?"
	result, err := t.db.Exec(
		command,
		task.Title,
		task.Description,
		task.Priority.String(),
		task.Status.String(),
		task.DueDate,
		task.CompletedAt,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("updateTask: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("updateTask: could not fetch rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("updateTask: no task found with id %s", task.ID)
	}

	return nil
}

func (t *taskRepository) Delete(id string) error {
	command := "DELETE FROM tasks WHERE id = ?"
	result, err := t.db.Exec(command, id)
	if err != nil {
		return fmt.Errorf("deleteTask: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteTask: could not fetch rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("deleteTask: no category found with id %s", id)
	}

	return nil
}
