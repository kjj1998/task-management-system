package repository

import "github.com/kjj1998/task-management-system/internal/models"

type UserRepository interface {
	GetById(id string) (models.DBUser, error)
	Create(user *models.DBUser) (string, error)
	Update(user *models.DBUser) error
	Delete(id string) error
}
