package user

import "github.com/kjj1998/task-management-system/internal/models"

type UserRepository interface {
	GetById(id string) (*models.DBUser, error)
	GetByEmail(email string) (*models.DBUser, error)
	Create(user *models.DBUser) (*models.DBUser, error)
	Update(user *models.DBUser) error
	Delete(id string) error
}
