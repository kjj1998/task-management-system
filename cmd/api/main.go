package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kjj1998/task-management-system/internal/logger"
	"github.com/kjj1998/task-management-system/internal/server"
)

func main() {
	env := os.Getenv("ENV")
	logger := logger.NewLogger(env)

	// Check environment
	if env == "dev" {
		logger.Info("Running in development environment")
	} else {
		logger.Info("Running in production environment")
	}

	server := server.NewTaskManagementSystemServer(logger)
	log.Fatal(http.ListenAndServe(":8080", server))
}
