package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kjj1998/task-management-system/internal/errors"
	"github.com/kjj1998/task-management-system/internal/models"
	"github.com/kjj1998/task-management-system/internal/services"
)

type TaskHandlers struct {
	taskService *services.TaskService
}

func NewTasksHandler(taskService *services.TaskService) *TaskHandlers {
	return &TaskHandlers{taskService: taskService}
}

func extractTaskID(path string) string {
	taskID := strings.TrimPrefix(path, "/tasks/")
	taskID = strings.TrimSuffix(taskID, "/")

	return taskID
}

func (h *TaskHandlers) HandleSingleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		taskID := extractTaskID(r.URL.Path)
		if taskID == "" {
			http.Error(w, "Task ID is required", http.StatusBadRequest)
			return
		}

		task, err := h.taskService.GetTask(w, taskID)

		if err != nil {
			errors.HandleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		}
	}
}

func (h *TaskHandlers) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID := r.URL.Query().Get("userId")
		if userID == "" {
			http.Error(w, "userId parameter is required", http.StatusBadRequest)
			return
		}

		tasks, err := h.taskService.GetTasksByUserID(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(tasks)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		}
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var task models.DBTask
		err = json.Unmarshal(body, &task)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		createdTask, err := h.taskService.CreateTask(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", fmt.Sprintf("/tasks/%s", createdTask.ID))
		w.WriteHeader(http.StatusCreated)

		response := models.CreateTaskResponse{
			ID:        createdTask.ID,
			Message:   "Task created successfully",
			CreatedAt: *createdTask.CreatedAt,
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		}
	case http.MethodDelete:

	}
}
