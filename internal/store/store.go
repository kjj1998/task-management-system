package store

import (
	"database/sql"
	"log/slog"

	"github.com/kjj1998/task-management-system/internal/errors"
	"github.com/kjj1998/task-management-system/internal/models"
	"github.com/kjj1998/task-management-system/internal/repository/category"
	"github.com/kjj1998/task-management-system/internal/repository/task"
	"github.com/kjj1998/task-management-system/internal/repository/user"
)

type DatabaseTaskStore struct {
	UserRepository     user.UserRepository
	CategoryRepository category.CategoryRepository
	TaskRepository     task.TaskRepository
}

func NewDatabaseTaskStore(db *sql.DB, errorHandler *errors.DatabaseErrorHandler, logger *slog.Logger) *DatabaseTaskStore {
	store := &DatabaseTaskStore{}

	store.UserRepository = user.NewUserRepository(db, errorHandler, logger)
	store.CategoryRepository = category.NewCategoryRepository(db, errorHandler, logger)
	store.TaskRepository = task.NewTaskRepository(db, errorHandler, logger)

	return store
}

func (s *DatabaseTaskStore) GetTask(taskID string) (*models.DBTask, error) {
	task, err := s.TaskRepository.GetById(taskID)
	if err != nil {
		return task, err
	}

	return task, nil
}
