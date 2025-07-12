package store

import (
	"database/sql"

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

func NewDatabaseTaskStore(db *sql.DB, errorHandler *errors.DatabaseErrorHandler) *DatabaseTaskStore {
	store := &DatabaseTaskStore{}

	store.UserRepository = user.NewUserRepository(db, errorHandler)
	store.CategoryRepository = category.NewCategoryRepository(db, errorHandler)
	store.TaskRepository = task.NewTaskRepository(db, errorHandler)

	return store
}

func (s *DatabaseTaskStore) GetTask(taskID string) (*models.DBTask, error) {
	task, err := s.TaskRepository.GetById(taskID)
	if err != nil {
		return task, err
	}

	return task, nil
}
