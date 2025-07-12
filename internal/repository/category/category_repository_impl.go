package category

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kjj1998/task-management-system/internal/errors"
	"github.com/kjj1998/task-management-system/internal/models"
)

const (
	createCategoryQuery     = "INSERT INTO categories (id, user_id, name, color) VALUES (?, ?, ?, ?)"
	getCategoryAfterCreate  = "SELECT id, created_at FROM categories WHERE id = ?"
	getAllCategoriesForUser = "SELECT id, user_id, name, color, created_at FROM categories WHERE user_id = ?"
	getCategoryByIDQuery    = "SELECT * FROM categories WHERE id = ?"
	updateCategoryQuery     = "UPDATE categories SET name = ?, color = ? WHERE id = ?"
)

type categoryRepository struct {
	db           *sql.DB
	errorHandler *errors.DatabaseErrorHandler
}

func NewCategoryRepository(db *sql.DB, errorHandler *errors.DatabaseErrorHandler) CategoryRepository {
	return &categoryRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

func (c *categoryRepository) scanDBCategory(rows any) (*models.DBCategory, error) {
	category := &models.DBCategory{}
	var err error
	switch r := rows.(type) {
	case *sql.Row:
		err = r.Scan(&category.ID, &category.UserID, &category.Name, &category.Color, &category.CreatedAt)
	case *sql.Rows:
		err = r.Scan(&category.ID, &category.UserID, &category.Name, &category.Color, &category.CreatedAt)
	default:
		return nil, fmt.Errorf("unsupported row type")
	}
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (c *categoryRepository) validateRowsAffected(result sql.Result, operation string, id string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.errorHandler.HandleDatabaseError(operation, err)
	}
	if rowsAffected == 0 {
		return c.errorHandler.HandleDatabaseError(operation, fmt.Errorf("no category found with id %s", id))
	}
	return nil
}

func (c *categoryRepository) Create(category *models.DBCategory) (*models.DBCategory, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, c.errorHandler.HandleDatabaseError("CreateCategory", err)
	}
	defer tx.Rollback()

	category_id := uuid.NewString()

	_, err = tx.Exec(createCategoryQuery, category_id, category.UserID, category.Name, category.Color)
	if err != nil {
		return nil, c.errorHandler.HandleDatabaseError("CreateCategory", err)
	}

	var createdCategory models.DBCategory
	err = tx.QueryRow(getCategoryAfterCreate, category_id).Scan(&createdCategory.ID, &createdCategory.CreatedAt)
	if err != nil {
		return nil, c.errorHandler.HandleDatabaseError("CreateCategory", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, c.errorHandler.HandleDatabaseError("CreateCategory", err)
	}

	return &createdCategory, nil
}

func (c *categoryRepository) GetAllForUser(user_id string) ([]models.DBCategory, error) {
	rows, err := c.db.Query(getAllCategoriesForUser, user_id)
	if err != nil {
		return nil, c.errorHandler.HandleDatabaseError("GetAllCategoriesForUser", err)
	}

	categories := make([]models.DBCategory, 0)
	for rows.Next() {
		task, err := c.scanDBCategory(rows)
		if err != nil {
			return nil, c.errorHandler.HandleDatabaseError("GetAllCategoriesForUser", err)
		}
		categories = append(categories, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, c.errorHandler.HandleDatabaseError("GetAllCategoriesForUser", err)
	}

	return categories, nil
}

func (c *categoryRepository) GetById(category_id string) (*models.DBCategory, error) {
	row := c.db.QueryRow(getCategoryByIDQuery, category_id)
	category, err := c.scanDBCategory(row)

	if err != nil {
		return nil, c.errorHandler.HandleDatabaseError("GetCategoryByID", err)
	}

	return category, nil
}

func (c *categoryRepository) Update(category *models.DBCategory) error {
	tx, err := c.db.Begin()
	if err != nil {
		return c.errorHandler.HandleDatabaseError("UpdateCategory", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(updateCategoryQuery, category.Name, category.Color, category.ID)
	if err != nil {
		return c.errorHandler.HandleDatabaseError("UpdateCategory", err)
	}

	if err := c.validateRowsAffected(result, "UpdateTask", category.ID); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return c.errorHandler.HandleDatabaseError("UpdateCategory", err)
	}

	return nil
}

func (c *categoryRepository) Delete(category_id string) error {
	command := "DELETE FROM categories WHERE id = ?"
	result, err := c.db.Exec(command, category_id)
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
