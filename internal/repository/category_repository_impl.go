package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kjj1998/task-management-system/internal/models"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(name string, color string, user_id string) (string, error) {
	command := "INSERT INTO categories (id, user_id, name, color) VALUES (?, ?, ?, ?)"
	id := uuid.NewString()
	_, err := r.db.Exec(command, id, user_id, name, color)
	if err != nil {
		return "", fmt.Errorf("createUser: %v", err)
	}

	return id, nil
}

func (r *categoryRepository) GetAllCategoriesForUser(user_id string) ([]models.DBCategory, error) {
	var categories []models.DBCategory

	query := "SELECT id, user_id, name, color, created_at FROM categories WHERE user_id = ?"
	rows, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, fmt.Errorf("categoriesByUser %q: %v", user_id, err)
	}

	for rows.Next() {
		var category models.DBCategory
		err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.Color, &category.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("categoriesByUser %q: %v", user_id, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("categoriesByUser %q: %v", user_id, err)
	}

	return categories, nil
}

func (r *categoryRepository) Update(name string, color string, id string, user_id string) error {
	fmt.Println("test")
	command := "UPDATE categories SET name = ?, color = ? WHERE id = ? AND user_id = ?"
	result, err := r.db.Exec(command, name, color, id, user_id)
	if err != nil {
		return fmt.Errorf("updateCategory: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	fmt.Println("rows affected: ", rowsAffected)
	if err != nil {
		return fmt.Errorf("updateCategory: could not fetch rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("updateCategory: no category found with id %s", id)
	}

	return nil
}

func (r *categoryRepository) Delete(id string, user_id string) error {
	command := "DELETE FROM categories WHERE id = ? AND user_id = ?"
	result, err := r.db.Exec(command, id, user_id)
	if err != nil {
		return fmt.Errorf("deleteCategory: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteCategory: could not fetch rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("deleteCategory: no category found with id %s", id)
	}

	return nil
}
