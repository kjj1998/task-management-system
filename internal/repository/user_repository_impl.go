package repository

import (
	"database/sql"
	"fmt"

	"github.com/kjj1998/task-management-system/internal/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.DbUser) (uint, error) {
	fmt.Println(user)
	query := "INSERT INTO users (email, password_hash, first_name, last_name) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, user.Email, user.PasswordHash, user.FirstName, user.LastName)
	if err != nil {
		return 0, fmt.Errorf("createUser: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("createUser: %v", err)
	}
	return uint(id), nil
}

func (r *userRepository) GetById(id uint) (user *models.DbUser, err error) {
	return &models.DbUser{}, nil
}

func (r *userRepository) Delete(id uint) error {
	return nil
}

func (r *userRepository) Update(user *models.DbUser) error {
	return nil
}
