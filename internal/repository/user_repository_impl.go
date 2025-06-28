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

func (r *userRepository) Create(user *models.DBUser) (uint, error) {
	command := "INSERT INTO users (email, password_hash, first_name, last_name) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(command, user.Email, user.PasswordHash, user.FirstName, user.LastName)
	if err != nil {
		return 0, fmt.Errorf("createUser: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"createUser: %v",
			err,
		)
	}
	return uint(id), nil
}

func (r *userRepository) GetById(id uint) (user models.DBUser, err error) {
	query := "SELECT * FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("usersById %d: no such user", id)
		}
		return user, fmt.Errorf("usersById %d: %v", id, err)
	}
	return user, nil
}

func (r *userRepository) Delete(id uint) error {
	command := "DELETE FROM users WHERE id = ?"
	result, err := r.db.Exec(command, id)
	if err != nil {
		return fmt.Errorf("deleteUser: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteUser: could not fetch rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("deleteUser: no user found with id %d", id)
	}

	return nil
}

func (r *userRepository) Update(user *models.DBUser) error {
	command := "UPDATE users SET email = ?, first_name = ?, last_name = ? WHERE id = ?"
	result, err := r.db.Exec(command, user.Email, user.FirstName, user.LastName, user.ID)
	if err != nil {
		return fmt.Errorf("updateUser: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("updateUser: could not fetch rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("updateUser: no user found with id %d", user.ID)
	}

	return nil
}
