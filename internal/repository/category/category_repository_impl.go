package category

import (
	"database/sql"
	"fmt"
	"log/slog"

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
	logger       *slog.Logger
}

func NewCategoryRepository(db *sql.DB, errorHandler *errors.DatabaseErrorHandler, logger *slog.Logger) CategoryRepository {
	return &categoryRepository{
		db:           db,
		errorHandler: errorHandler,
		logger:       logger,
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
	c.logger.Debug("creating category", slog.String("user_id", category.UserID))

	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Error("failed to create category", slog.String("error", err.Error()))
		return nil, c.errorHandler.HandleDatabaseError("CreateCategory", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			c.logger.Warn("failed to rollback transaction", slog.String("error", rollbackErr.Error()))
		}
	}()

	category_id := uuid.NewString()

	_, err = tx.Exec(createCategoryQuery, category_id, category.UserID, category.Name, category.Color)
	if err != nil {
		c.logger.Error("failed to create category", slog.String("error", err.Error()))
		return nil, c.errorHandler.HandleDatabaseError("CreateCategory", err)
	}

	var createdCategory models.DBCategory
	err = tx.QueryRow(getCategoryAfterCreate, category_id).Scan(&createdCategory.ID, &createdCategory.CreatedAt)
	if err != nil {
		c.logger.Error("failed to create category", slog.String("error", err.Error()))
		return nil, c.errorHandler.HandleDatabaseError("CreateCategory", err)
	}

	err = tx.Commit()
	if err != nil {
		c.logger.Error("failed to create category", slog.String("error", err.Error()))
		return nil, c.errorHandler.HandleDatabaseError("CreateCategory", err)
	}

	c.logger.Info("category created", slog.String("category_id", createdCategory.ID))
	return &createdCategory, nil
}

func (c *categoryRepository) GetAllForUser(user_id string) ([]models.DBCategory, error) {
	c.logger.Debug("fetching categories", slog.String("user_id", user_id))

	rows, err := c.db.Query(getAllCategoriesForUser, user_id)
	if err != nil {
		c.logger.Error("failed to fetch categories", slog.String("error", err.Error()))
		return nil, c.errorHandler.HandleDatabaseError("GetAllCategoriesForUser", err)
	}

	categories := make([]models.DBCategory, 0)
	for rows.Next() {
		category, err := c.scanDBCategory(rows)
		if err != nil {
			c.logger.Error("failed to scan category", slog.String("error", err.Error()))
			return nil, c.errorHandler.HandleDatabaseError("GetAllCategoriesForUser", err)
		}
		categories = append(categories, *category)
	}

	if err := rows.Err(); err != nil {
		c.logger.Error("error reading categories", slog.String("error", err.Error()))
		return nil, c.errorHandler.HandleDatabaseError("GetAllCategoriesForUser", err)
	}

	c.logger.Debug("categories retrieved", slog.String("user_id", user_id), slog.Int("count", len(categories)))
	return categories, nil
}

func (c *categoryRepository) GetById(category_id string) (*models.DBCategory, error) {
	c.logger.Debug("fetching category", slog.String("category_id", category_id))

	row := c.db.QueryRow(getCategoryByIDQuery, category_id)
	category, err := c.scanDBCategory(row)

	if err != nil {
		c.logger.Error("failed to fetch category", slog.String("error", err.Error()))
		return nil, c.errorHandler.HandleDatabaseError("GetCategoryByID", err)
	}

	c.logger.Info("category retrieved", slog.String("category_id", category.ID))
	return category, nil
}

func (c *categoryRepository) Update(category *models.DBCategory) error {
	c.logger.Debug("updating category", slog.String("category_id", category.ID))

	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Error("failed to update category", slog.String("error", err.Error()))
		return c.errorHandler.HandleDatabaseError("UpdateCategory", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			c.logger.Warn("failed to rollback transaction", slog.String("error", rollbackErr.Error()))
		}
	}()

	result, err := tx.Exec(updateCategoryQuery, category.Name, category.Color, category.ID)
	if err != nil {
		c.logger.Error("failed to update category", slog.String("error", err.Error()))
		return c.errorHandler.HandleDatabaseError("UpdateCategory", err)
	}

	if err := c.validateRowsAffected(result, "UpdateCategory", category.ID); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		c.logger.Error("failed to update category", slog.String("error", err.Error()))
		return c.errorHandler.HandleDatabaseError("UpdateCategory", err)
	}

	c.logger.Info("category updated", slog.String("category_id", category.ID))
	return nil
}

func (c *categoryRepository) Delete(category_id string) error {
	c.logger.Debug("deleting category", slog.String("category_id", category_id))

	command := "DELETE FROM categories WHERE id = ?"
	result, err := c.db.Exec(command, category_id)
	if err != nil {
		c.logger.Error("failed to delete category", slog.String("error", err.Error()))
		return c.errorHandler.HandleDatabaseError("DeleteCategory", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.logger.Error("failed to check deletion result", slog.String("error", err.Error()))
		return c.errorHandler.HandleDatabaseError("DeleteCategory", err)
	}
	if rowsAffected == 0 {
		c.logger.Warn("category not found for deletion", slog.String("category_id", category_id))
		return c.errorHandler.HandleDatabaseError("DeleteCategory", err)
	}

	c.logger.Info("category deleted", slog.String("category_id", category_id))
	return nil
}
