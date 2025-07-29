package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kjj1998/task-management-system/internal/errors"
	"github.com/kjj1998/task-management-system/internal/models"
	"github.com/kjj1998/task-management-system/internal/services"
)

type TaskHandlers struct {
	taskService *services.TaskService
	logger      *slog.Logger
}

func NewTasksHandler(taskService *services.TaskService, logger *slog.Logger) *TaskHandlers {
	return &TaskHandlers{taskService: taskService, logger: logger}
}

func extractTaskID(path string) string {
	taskID := strings.TrimPrefix(path, "/tasks/")
	taskID = strings.TrimSuffix(taskID, "/")

	return taskID
}

func (h *TaskHandlers) HandleSingleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTaskByID(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandlers) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := extractTaskID(r.URL.Path)
	if taskID == "" {
		validationError := errors.NewBadRequestError("Task ID is required", nil)
		errors.HandleError(w, validationError, h.logger)
		return
	}

	task, err := h.taskService.GetTask(w, taskID)
	if err != nil {
		errors.HandleError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := models.NewSuccessResponse("Task retrieved successfully", task)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		errors.HandleError(w, err, h.logger)
	}
}

func (h *TaskHandlers) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTasks(w, r)
	case http.MethodPost:
		h.CreateTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		validationError := errors.NewBadRequestError("User ID parameter is required", nil)
		errors.HandleError(w, validationError, h.logger)
		return
	}

	tasks, err := h.taskService.GetTasksByUserID(userID)
	if err != nil {
		errors.HandleError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := models.NewSuccessResponse("Tasks retrieved successfully", tasks)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		errors.HandleError(w, err, h.logger)
	}
}

func (h *TaskHandlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		parsingError := errors.NewBadRequestError("Error reading request body", nil)
		errors.HandleError(w, parsingError, h.logger)
		return
	}
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			h.logger.Warn("failed to close request body", slog.String("error", closeErr.Error()))
		}
	}()

	var task models.DBTask
	err = json.Unmarshal(body, &task)
	if err != nil {
		parsingError := errors.NewBadRequestError("Error parsing json body", nil)
		errors.HandleError(w, parsingError, h.logger)
		return
	}

	createdTask, err := h.taskService.CreateTask(task)
	if err != nil {
		errors.HandleError(w, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/tasks/%s", createdTask.ID))
	w.WriteHeader(http.StatusCreated)

	response := models.NewSuccessResponse("Task created successfully", createdTask)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		errors.HandleError(w, err, h.logger)
	}
}
