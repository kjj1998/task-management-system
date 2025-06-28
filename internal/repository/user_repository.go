package repository

import "github.com/kjj1998/task-management-system/internal/models"

type UserRepository interface {
	GetById(id uint) (*models.DbUser, error)
	Create(user *models.DbUser) (uint, error)
	Update(user *models.DbUser) error
	Delete(id uint) error
}
