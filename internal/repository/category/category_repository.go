package category

import "github.com/kjj1998/task-management-system/internal/models"

type CategoryRepository interface {
	GetAllForUser(user_id string) ([]models.DBCategory, error)
	GetById(category_id string) (*models.DBCategory, error)
	Create(*models.DBCategory) (*models.DBCategory, error)
	Update(*models.DBCategory) error
	Delete(category_id string) error
}
