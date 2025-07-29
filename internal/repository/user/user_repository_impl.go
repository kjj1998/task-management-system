package user

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/kjj1998/task-management-system/internal/errors"
	"github.com/kjj1998/task-management-system/internal/models"
)

const (
	createUserQuery         = "INSERT INTO users (id, email, password_hash, first_name, last_name) VALUES (?, ?, ?, ?, ?)"
	getUserAfterCreateQuery = "SELECT id, created_at FROM users WHERE id = ?"
	getUserByIDQuery        = "SELECT * FROM users WHERE id = ?"
	getUserByEmail          = "SELECT * FROM users WHERE email = ?"
	deleteUserQuery         = "DELETE FROM users WHERE id = ?"
	updateUserQuery         = "UPDATE users SET email = ?, first_name = ?, last_name = ? WHERE id = ?"
)

type userRepository struct {
	db           *sql.DB
	errorHandler *errors.DatabaseErrorHandler
	logger       *slog.Logger
}

func NewUserRepository(db *sql.DB, errorHandler *errors.DatabaseErrorHandler, logger *slog.Logger) UserRepository {
	return &userRepository{
		db:           db,
		errorHandler: errorHandler,
		logger:       logger,
	}
}

func (u *userRepository) scanDBUser(rows any) (*models.DBUser, error) {
	user := &models.DBUser{}
	var err error
	switch r := rows.(type) {
	case *sql.Row:
		err = r.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	case *sql.Rows:
		err = r.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	default:
		return nil, fmt.Errorf("unsupported row type")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) validateRowsAffected(result sql.Result, operation string, id string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return u.errorHandler.HandleDatabaseError(operation, err)
	}
	if rowsAffected == 0 {
		return u.errorHandler.HandleDatabaseError(operation, fmt.Errorf("no user found with id %s", id))
	}
	return nil
}

func (u *userRepository) Create(user *models.DBUser) (*models.DBUser, error) {
	u.logger.Debug("creating user", slog.String("email", user.Email))

	tx, err := u.db.Begin()
	if err != nil {
		return nil, u.errorHandler.HandleDatabaseError("CreateUser", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			u.logger.Warn("failed to rollback transaction", slog.String("error", rollbackErr.Error()))
		}
	}()

	userID := uuid.NewString()

	_, err = tx.Exec(createUserQuery, userID, user.Email, user.PasswordHash, user.FirstName, user.LastName)
	if err != nil {
		return nil, u.errorHandler.HandleDatabaseError("CreateUser", err)
	}

	var createdUser models.DBUser
	err = tx.QueryRow(getUserAfterCreateQuery, userID).Scan(&createdUser.ID, &createdUser.CreatedAt)
	if err != nil {
		return nil, u.errorHandler.HandleDatabaseError("CreateUser", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, u.errorHandler.HandleDatabaseError("CreateUser", err)
	}

	u.logger.Info("created user", slog.String("user_id", createdUser.ID))
	return &createdUser, nil
}

func (u *userRepository) GetById(id string) (*models.DBUser, error) {
	u.logger.Debug("get user by id", slog.String("user_id", id))

	row := u.db.QueryRow(getUserByIDQuery, id)
	user, err := u.scanDBUser(row)
	if err != nil {
		return nil, u.errorHandler.HandleDatabaseError("GetUserByID", err)
	}

	u.logger.Info("got user", slog.String("user_id", id))
	return user, nil
}

func (u *userRepository) GetByEmail(email string) (*models.DBUser, error) {
	u.logger.Debug("get user by email", slog.String("email", email))

	row := u.db.QueryRow(getUserByEmail, email)
	user, err := u.scanDBUser(row)
	if err != nil {
		return nil, u.errorHandler.HandleDatabaseError("GetUserByEmail", err)
	}

	u.logger.Info("got user by email", slog.String("email", email))
	return user, nil
}

func (u *userRepository) Delete(id string) error {
	u.logger.Debug("delete user", slog.String("user_id", id))

	tx, err := u.db.Begin()
	if err != nil {
		return u.errorHandler.HandleDatabaseError("DeleteUser", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			u.logger.Warn("failed to rollback transaction", slog.String("error", rollbackErr.Error()))
		}
	}()

	result, err := tx.Exec(deleteUserQuery, id)
	if err != nil {
		return u.errorHandler.HandleDatabaseError("DeleteUser", err)
	}

	if err := u.validateRowsAffected(result, "DeleteUser", id); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return u.errorHandler.HandleDatabaseError("DeleteUser", err)
	}

	u.logger.Info("deleted user", slog.String("user_id", id))
	return nil
}

func (u *userRepository) Update(user *models.DBUser) error {
	u.logger.Debug("update user", slog.String("user_id", user.ID))

	tx, err := u.db.Begin()
	if err != nil {
		return u.errorHandler.HandleDatabaseError("UpdateUser", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			u.logger.Warn("failed to rollback transaction", slog.String("error", rollbackErr.Error()))
		}
	}()

	result, err := tx.Exec(
		updateUserQuery,
		user.Email,
		user.FirstName,
		user.LastName,
		user.ID,
	)

	if err != nil {
		return u.errorHandler.HandleDatabaseError("UpdateUser", err)
	}

	if err := u.validateRowsAffected(result, "UpdateUser", user.ID); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return u.errorHandler.HandleDatabaseError("UpdateUser", err)
	}

	u.logger.Info("updated user", slog.String("user_id", user.ID))
	return nil
}
