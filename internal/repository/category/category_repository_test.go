package category_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/kjj1998/task-management-system/internal/database"
	"github.com/kjj1998/task-management-system/internal/models"
	"github.com/kjj1998/task-management-system/internal/repository/category"
	"github.com/kjj1998/task-management-system/internal/repository/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CategoryRepoTestSuite struct {
	suite.Suite
	mySQLContainer *testutils.MySQLContainer
	ctx            context.Context
	repository     category.CategoryRepository
}

func (suite *CategoryRepoTestSuite) SetupSuite() {
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
	categoryRepository := category.NewCategoryRepository(db)
	suite.repository = categoryRepository
}

func (suite *CategoryRepoTestSuite) TearDownSuite() {
	if err := suite.mySQLContainer.Container.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating mysql container: %s", err)
	}
}

func (suite *CategoryRepoTestSuite) TestCreateCategory() {
	t := suite.T()

	category := &models.DBCategory{
		UserID: "1244ABC",
		Name:   "urgent",
		Color:  "#ff0000",
	}
	category_id, err := suite.repository.Create(category)
	assert.NoError(t, err)
	assert.NotNil(t, category_id)
}

func (suite *CategoryRepoTestSuite) TestGetAllCategoriesForUser() {
	t := suite.T()

	categories, err := suite.repository.GetAllForUser("1244ABC")

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Len(t, categories, 2)
}

func (suite *CategoryRepoTestSuite) TestUpdateCategory() {
	t := suite.T()

	category, err := suite.repository.GetById("2345SDSXAS")
	assert.NoError(t, err)
	assert.NotNil(t, category)

	updatedCategory := &models.DBCategory{
		ID:     category.ID,
		UserID: category.UserID,
		Name:   "Do by today",
		Color:  "#ffff00",
	}
	suite.repository.Update(updatedCategory)

	category, err = suite.repository.GetById("2345SDSXAS")
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "Do by today", category.Name)
	assert.Equal(t, "#ffff00", category.Color)
}

func (suite *CategoryRepoTestSuite) TestDeleteCategory() {
	t := suite.T()

	category := models.DBCategory{
		UserID: "1244ABC",
		Name:   "super urgent",
		Color:  "#0000ff",
	}
	category_id, err := suite.repository.Create(&category)
	assert.NoError(t, err)
	assert.NotNil(t, category_id)

	category, err = suite.repository.GetById(category_id)
	assert.NoError(t, err)
	assert.NotNil(t, category)

	err = suite.repository.Delete(category.ID)
	assert.NoError(t, err)

	_, err = suite.repository.GetById(category_id)
	expectedErrorMessage := fmt.Sprintf("categoryById %s: no such category", category_id)
	assert.Contains(t, err.Error(), expectedErrorMessage)
}

func TestCategoryRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepoTestSuite))
}
