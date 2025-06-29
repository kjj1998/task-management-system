package repository

import "github.com/kjj1998/task-management-system/internal/models"

type TaskRepository interface {
	Create(task *models.DBTask) (string, error)
	GetAllTasksForUser(user_id string) ([]models.DBTask, error)
	Update(task *models.DBTask) error
	Delete(task_id string) error
}
