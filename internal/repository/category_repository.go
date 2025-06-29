package repository

import "github.com/kjj1998/task-management-system/internal/models"

type CategoryRepository interface {
	GetAllCategoriesForUser(user_id string) ([]models.DBCategory, error)
	Create(name string, color string, user_id string) (string, error)
	Update(name string, color string, id string, user_id string) error
	Delete(id string, user_id string) error
}
