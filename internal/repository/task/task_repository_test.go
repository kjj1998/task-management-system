package task_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/kjj1998/task-management-system/internal/database"
	"github.com/kjj1998/task-management-system/internal/errors"
	"github.com/kjj1998/task-management-system/internal/logger"
	"github.com/kjj1998/task-management-system/internal/models"
	"github.com/kjj1998/task-management-system/internal/repository/task"
	"github.com/kjj1998/task-management-system/internal/repository/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TaskRepoTestSuite struct {
	suite.Suite
	mySQLContainer *testutils.MySQLContainer
	ctx            context.Context
	repository     task.TaskRepository
}

func (suite *TaskRepoTestSuite) SetupSuite() {
	logger := logger.NewLogger("test")
	suite.ctx = context.Background()

	mySQLContainer, err := testutils.CreateMySQLContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.mySQLContainer = mySQLContainer
	host, _ := mySQLContainer.Container.Host(suite.ctx)
	port, _ := mySQLContainer.Container.MappedPort(suite.ctx, "3306")

	database.Connect("testuser", "testpass", host, port.Port(), "taskapi", logger)
	db := database.GetDb()
	dbErrorHandler := errors.NewDatabaseErrorHandler()
	taskRepository := task.NewTaskRepository(db, dbErrorHandler, logger)
	suite.repository = taskRepository
}

func (suite *TaskRepoTestSuite) TearDownSuite() {
	if err := suite.mySQLContainer.Container.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating mysql container: %s", err)
	}
}

func (suite *TaskRepoTestSuite) TestTaskRepositoryOperations() {
	t := suite.T()

	t.Run("CreateTask", func(t *testing.T) {
		dueDate := time.Date(2025, time.July, 3, 22, 18, 0, 0, time.UTC)
		task := &models.DBTask{
			UserID:      "1244ABC",
			CategoryID:  "2345SDSXAS",
			Title:       "Collect Parcel",
			Description: "Collect parcel from the delivery point",
			Priority:    models.Medium,
			Status:      models.Pending,
			DueDate:     &dueDate,
		}

		createdTask, err := suite.repository.Create(task)
		assert.NoError(t, err)
		assert.NotNil(t, createdTask)

		if t.Failed() {
			t.Fatal("CreateTask failed, stopping sequential execution")
		}
	})

	t.Run("GetAllTasksForUser", func(t *testing.T) {
		tasks, err := suite.repository.GetAllForUser("1244ABC")

		assert.NoError(t, err)
		assert.NotNil(t, tasks)
		assert.Len(t, tasks, 2)
	})

	t.Run("UpdateTask", func(t *testing.T) {
		dueDate := time.Date(2025, time.July, 3, 22, 18, 0, 0, time.UTC)
		completedTime := time.Date(2025, time.July, 4, 22, 18, 0, 0, time.UTC)

		task := &models.DBTask{
			ID:          "DSFDS23423",
			UserID:      "1244ABC",
			CategoryID:  "2345SDSXAS",
			Title:       "Collect Parcel",
			Description: "Collect parcel from the delivery point",
			Priority:    models.Medium,
			Status:      models.Completed,
			DueDate:     &dueDate,
			CompletedAt: &completedTime,
		}

		err := suite.repository.Update(task)
		assert.NoError(t, err)

		updatedTask, err := suite.repository.GetById("DSFDS23423")
		assert.NoError(t, err)
		assert.NotNil(t, updatedTask)
		assert.Equal(t, "Collect Parcel", updatedTask.Title)
		assert.Equal(t, "Collect parcel from the delivery point", updatedTask.Description)
		assert.Equal(t, models.Completed, updatedTask.Status)
		assert.Equal(t, &dueDate, updatedTask.DueDate)
		assert.Equal(t, &completedTime, updatedTask.CompletedAt)
	})

	t.Run("DeleteTask", func(t *testing.T) {
		err := suite.repository.Delete("DSFDS23423")
		assert.NoError(t, err)

		_, err = suite.repository.GetById("DSFDS23423")
		expectedErrorMessage := "Resource not found, sql: no rows in result set"
		assert.Contains(t, err.Error(), expectedErrorMessage)
	})
}

func TestCategoryRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepoTestSuite))
}
