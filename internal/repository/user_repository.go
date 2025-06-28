package repository

import "github.com/kjj1998/task-management-system/internal/models"

type UserRepository interface {
	GetById(id uint) (models.DBUser, error)
	Create(user *models.DBUser) (uint, error)
	Update(user *models.DBUser) error
	Delete(id uint) error
}
