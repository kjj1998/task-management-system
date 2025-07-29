package services

import (
	"net/http"

	"github.com/kjj1998/task-management-system/internal/models"
	"github.com/kjj1998/task-management-system/internal/store"
)

type TaskService struct {
	taskStore *store.DatabaseTaskStore
}

func NewTaskService(taskStore *store.DatabaseTaskStore) *TaskService {
	return &TaskService{
		taskStore: taskStore,
	}
}

func (s *TaskService) GetTask(w http.ResponseWriter, task_id string) (*models.DBTask, error) {
	task, err := s.taskStore.TaskRepository.GetById(task_id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTasksByUserID(user_id string) ([]models.DBTask, error) {
	tasks, err := s.taskStore.TaskRepository.GetAllForUser(user_id)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) CreateTask(task models.DBTask) (*models.DBTask, error) {
	createdTask, err := s.taskStore.TaskRepository.Create(&task)
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}
