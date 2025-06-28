package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() error {
	// Capture connection properties
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASS")
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	cfg.DBName = os.Getenv("DB")
	cfg.Params = map[string]string{"parseTime": "true"}

	// Get a database handle
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return fmt.Errorf("failed to ping db: %w", pingErr)
	}
	fmt.Println("Database Connected!")
	DB = db

	return nil
}

func Close() error {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			return fmt.Errorf("failed to close DB connection: %w", err)
		}
	}

	return nil
}

func GetDb() *sql.DB {
	return DB
}
