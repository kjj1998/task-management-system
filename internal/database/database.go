package database

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect(user, password, host, port, dbName string) error {
	// Capture connection properties
	cfg := mysql.NewConfig()
	cfg.User = user
	cfg.Passwd = password
	cfg.Net = "tcp"
	cfg.Addr = host + ":" + port
	cfg.DBName = dbName
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
