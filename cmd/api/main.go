package main

import (
	"log"
	"net/http"

	"github.com/kjj1998/task-management-system/internal/config"
	"github.com/kjj1998/task-management-system/internal/logger"
	"github.com/kjj1998/task-management-system/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger := logger.NewLogger(cfg.Environment)

	if cfg.IsDevelopment() {
		logger.Info("Running in development environment")
	} else {
		logger.Info("Running in production environment")
	}

	server := server.NewTaskManagementSystemServer(cfg, logger)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, server))
}
