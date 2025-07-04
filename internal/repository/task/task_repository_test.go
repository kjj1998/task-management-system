package task_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/kjj1998/task-management-system/internal/database"
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
	suite.ctx = context.Background()

	mySQLContainer, err := testutils.CreateMySQLContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.mySQLContainer = mySQLContainer
	host, _ := mySQLContainer.Container.Host(suite.ctx)
	port, _ := mySQLContainer.Container.MappedPort(suite.ctx, "3306")

	database.Connect("testuser", "testpass", host, port.Port(), "taskapi")
	db := database.GetDb()
	taskRepository := task.NewTaskRepository(db)
	suite.repository = taskRepository
}

func (suite *TaskRepoTestSuite) TearDownSuite() {
	if err := suite.mySQLContainer.Container.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating mysql container: %s", err)
	}
}

func (suite *TaskRepoTestSuite) TestCreateTask() {
	t := suite.T()

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

	task_id, err := suite.repository.Create(task)
	assert.NoError(t, err)
	assert.NotNil(t, task_id)
}

func (suite *TaskRepoTestSuite) TestGetAllCategoriesForUser() {
	t := suite.T()

	categories, err := suite.repository.GetAllForUser("1244ABC")

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Len(t, categories, 2)
}

func (suite *TaskRepoTestSuite) TestUpdateTask() {
	t := suite.T()

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
}

func (suite *TaskRepoTestSuite) TestDeleteTask() {
	t := suite.T()

	dueDate := time.Date(2025, time.July, 3, 22, 18, 0, 0, time.UTC)
	task := &models.DBTask{
		UserID:      "1244ABC",
		CategoryID:  "2345SDSXAS",
		Title:       "Take out the trash",
		Description: "take the trash out",
		Priority:    models.Medium,
		Status:      models.Pending,
		DueDate:     &dueDate,
	}

	task_id, err := suite.repository.Create(task)
	assert.NoError(t, err)
	assert.NotNil(t, task_id)

	err = suite.repository.Delete(task_id)
	assert.NoError(t, err)

	_, err = suite.repository.GetById(task_id)
	expectedErrorMessage := fmt.Sprintf("taskById %s: no such task", task_id)
	assert.Contains(t, err.Error(), expectedErrorMessage)
}

func TestCategoryRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepoTestSuite))
}
