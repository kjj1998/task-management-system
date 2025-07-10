package task

import "github.com/kjj1998/task-management-system/internal/models"

type TaskRepository interface {
	Create(task *models.DBTask) (*models.DBTask, error)
	GetAllForUser(user_id string) ([]models.DBTask, error)
	GetById(task_id string) (*models.DBTask, error)
	Update(task *models.DBTask) error
	Delete(task_id string) error
}
