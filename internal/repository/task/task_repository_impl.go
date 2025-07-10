package task

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kjj1998/task-management-system/internal/errors"
	"github.com/kjj1998/task-management-system/internal/models"
)

const (
	createTaskQuery    = "INSERT INTO tasks (id, user_id, category_id, title, description, priority, status, due_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	getTaskByIDQuery   = "SELECT * FROM tasks WHERE id = ?"
	getTaskAfterCreate = "SELECT id, created_at FROM tasks WHERE id = ?"
	getAllTasksForUser = "SELECT id, user_id, category_id, title, description, priority, status, due_date, completed_at, created_at, updated_at FROM tasks WHERE user_id = ?"
	updateTaskQuery    = "UPDATE tasks SET title = ?, description = ?, priority = ?, status = ?, due_date = ?, completed_at = ?, updated_at = ? WHERE id = ?"
	deleteTaskQuery    = "DELETE FROM tasks WHERE id = ?"
)

type taskRepository struct {
	db           *sql.DB
	errorHandler *errors.DatabaseErrorHandler
}

func NewTaskRepository(db *sql.DB, errorHandler *errors.DatabaseErrorHandler) TaskRepository {
	return &taskRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

func (t *taskRepository) scanDBTask(rows any) (*models.DBTask, error) {
	task := &models.DBTask{}
	var err error
	switch r := rows.(type) {
	case *sql.Row:
		err = r.Scan(&task.ID, &task.UserID, &task.CategoryID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.DueDate, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt)
	case *sql.Rows:
		err = r.Scan(&task.ID, &task.UserID, &task.CategoryID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.DueDate, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt)
	default:
		return nil, fmt.Errorf("unsupported row type")
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *taskRepository) validateRowsAffected(result sql.Result, operation string, id string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return t.errorHandler.HandleDatabaseError(operation, err)
	}
	if rowsAffected == 0 {
		return t.errorHandler.HandleDatabaseError(operation, fmt.Errorf("no task found with id %s", id))
	}
	return nil
}

func (t *taskRepository) Create(task *models.DBTask) (*models.DBTask, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, t.errorHandler.HandleDatabaseError("CreateTask", err)
	}
	defer tx.Rollback()

	task_id := uuid.NewString()

	_, err = tx.Exec(createTaskQuery, task_id, task.UserID, task.CategoryID, task.Title, task.Description, task.Priority, task.Status, task.DueDate)
	if err != nil {
		return nil, t.errorHandler.HandleDatabaseError("CreateTask", err)
	}

	var createdTask models.DBTask
	err = tx.QueryRow(getTaskAfterCreate, task_id).Scan(&createdTask.ID, &createdTask.CreatedAt)
	if err != nil {
		return nil, t.errorHandler.HandleDatabaseError("CreateTask", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, t.errorHandler.HandleDatabaseError("CreateTask", err)
	}

	return &createdTask, nil
}

func (t *taskRepository) GetAllForUser(user_id string) ([]models.DBTask, error) {
	rows, err := t.db.Query(getAllTasksForUser, user_id)
	if err != nil {
		return nil, t.errorHandler.HandleDatabaseError("GetAllTasksForUser", err)
	}

	tasks := make([]models.DBTask, 0)
	for rows.Next() {
		task, err := t.scanDBTask(rows)
		if err != nil {
			return nil, t.errorHandler.HandleDatabaseError("GetAllTasksForUser", err)
		}
		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, t.errorHandler.HandleDatabaseError("GetAllTasksForUser", err)
	}

	return tasks, nil
}

func (t *taskRepository) GetById(task_id string) (*models.DBTask, error) {
	row := t.db.QueryRow(getTaskByIDQuery, task_id)
	task, err := t.scanDBTask(row)
	if err != nil {
		return nil, t.errorHandler.HandleDatabaseError("GetTaskByID", err)
	}

	return task, nil
}

func (t *taskRepository) Update(task *models.DBTask) error {
	tx, err := t.db.Begin()
	if err != nil {
		return t.errorHandler.HandleDatabaseError("UpdateTask", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		updateTaskQuery,
		task.Title,
		task.Description,
		task.Priority,
		task.Status,
		task.DueDate,
		task.CompletedAt,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return t.errorHandler.HandleDatabaseError("UpdateTask", err)
	}

	if err := t.validateRowsAffected(result, "UpdateTask", task.ID); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return t.errorHandler.HandleDatabaseError("UpdateTask", err)
	}

	return nil
}

func (t *taskRepository) Delete(id string) error {
	tx, err := t.db.Begin()
	if err != nil {
		return t.errorHandler.HandleDatabaseError("DeleteTask", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(deleteTaskQuery, id)
	if err != nil {
		return t.errorHandler.HandleDatabaseError("DeleteTask", err)
	}

	if err := t.validateRowsAffected(result, "DeleteTask", id); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return t.errorHandler.HandleDatabaseError("DeleteTask", err)
	}

	return nil
}
