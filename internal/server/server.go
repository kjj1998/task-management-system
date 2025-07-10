package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kjj1998/task-management-system/internal/database"
	"github.com/kjj1998/task-management-system/internal/errors"
	"github.com/kjj1998/task-management-system/internal/handlers"
	"github.com/kjj1998/task-management-system/internal/services"
	"github.com/kjj1998/task-management-system/internal/store"
)

type TaskManagementSystemStore struct{}

type TaskManagementSystemServer struct {
	http.Handler
}

func NewTaskManagementSystemServer() *TaskManagementSystemServer {
	err := database.Connect(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
	}
	db := database.GetDb()
	dbErrorHandler := errors.NewDatabaseErrorHandler()

	store := store.NewDatabaseTaskStore(db, dbErrorHandler)
	taskService := services.NewTaskService(store)
	taskHandler := handlers.NewTasksHandler(taskService)

	t := new(TaskManagementSystemServer)

	router := http.NewServeMux()
	router.Handle("/tasks/", http.HandlerFunc(taskHandler.HandleSingleTask))
	router.Handle("/tasks", http.HandlerFunc(taskHandler.HandleTasks))
	router.Handle("/healthcheck", http.HandlerFunc(t.healthcheckHandler))

	t.Handler = http.StripPrefix("/api", router)

	return t
}

func (t *TaskManagementSystemServer) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]string{"healthcheck": "online"}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(health)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
