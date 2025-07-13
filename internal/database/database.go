package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect(user, password, host, port, dbName string, logger *slog.Logger) error {
	// Capture connection properties
	cfg := mysql.NewConfig()
	cfg.User = user
	cfg.Passwd = password
	cfg.Net = "tcp"
	cfg.Addr = host + ":" + port
	cfg.DBName = dbName
	cfg.Params = map[string]string{"parseTime": "true"}

	logger.Info("attempting database connection",
		slog.String("host", host),
		slog.String("port", port),
		slog.String("database", dbName),
		slog.String("user", user),
	)

	// Get a database handle
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Error("failed to open database connection",
			slog.String("error", err.Error()),
			slog.String("host", host),
			slog.String("port", port),
		)
		return fmt.Errorf("failed to open db: %w", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		logger.Error("failed to ping database",
			slog.String("error", pingErr.Error()),
			slog.String("host", host),
			slog.String("port", port),
		)
		return fmt.Errorf("failed to ping db: %w", pingErr)
	}

	logger.Info("database connection established successfully")
	DB = db

	return nil
}

func Close(logger *slog.Logger) error {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			logger.Error("failed to close database connection", slog.String("error", err.Error()))
			return fmt.Errorf("failed to close DB connection: %w", err)
		}
		logger.Info("database connection closed")
	}

	return nil
}

func GetDb() *sql.DB {
	return DB
}
