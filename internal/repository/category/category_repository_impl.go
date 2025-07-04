package category

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

func (r *categoryRepository) Create(category *models.DBCategory) (string, error) {
	command := "INSERT INTO categories (id, user_id, name, color) VALUES (?, ?, ?, ?)"

	category_id := uuid.NewString()
	_, err := r.db.Exec(command, category_id, category.UserID, category.Name, category.Color)
	if err != nil {
		return "", fmt.Errorf("createUser: %v", err)
	}

	return category_id, nil
}

func (r *categoryRepository) GetAllForUser(user_id string) ([]models.DBCategory, error) {
	query := "SELECT id, user_id, name, color, created_at FROM categories WHERE user_id = ?"
	rows, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, fmt.Errorf("categoriesByUser %q: %v", user_id, err)
	}

	var categories []models.DBCategory
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

func (r *categoryRepository) GetById(category_id string) (category models.DBCategory, err error) {
	query := "SELECT * FROM categories WHERE id = ?"
	row := r.db.QueryRow(query, category_id)

	err = row.Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.Color,
		&category.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return category, fmt.Errorf("categoryById %s: no such category", category_id)
		}
		return category, fmt.Errorf("categoryById %s: %v", category_id, err)
	}

	return category, nil
}

func (r *categoryRepository) Update(category *models.DBCategory) error {
	command := "UPDATE categories SET name = ?, color = ? WHERE id = ?"
	result, err := r.db.Exec(command, category.Name, category.Color, category.ID)
	if err != nil {
		return fmt.Errorf("updateCategory: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("updateCategory: could not fetch rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("updateCategory: no category found with id %s", category.ID)
	}

	return nil
}

func (r *categoryRepository) Delete(category_id string) error {
	command := "DELETE FROM categories WHERE id = ?"
	result, err := r.db.Exec(command, category_id)
	if err != nil {
		return fmt.Errorf("deleteCategory: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("deleteCategory: could not fetch rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("deleteCategory: no category found with id %s", category_id)
	}

	return nil
}
