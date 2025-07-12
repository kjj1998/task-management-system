package category_test

import (
	"context"
	"log"
	"testing"

	"github.com/kjj1998/task-management-system/internal/database"
	"github.com/kjj1998/task-management-system/internal/errors"
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
	dbErrorHandler := errors.NewDatabaseErrorHandler()
	categoryRepository := category.NewCategoryRepository(db, dbErrorHandler)
	suite.repository = categoryRepository
}

func (suite *CategoryRepoTestSuite) TearDownSuite() {
	if err := suite.mySQLContainer.Container.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating mysql container: %s", err)
	}
}

func (suite *CategoryRepoTestSuite) TestCategoryRepositoryOperations() {
	t := suite.T()

	t.Run("CreateCategory", func(t *testing.T) {
		category := &models.DBCategory{
			UserID: "1244ABC",
			Name:   "urgent",
			Color:  "#ff0000",
		}

		createdCategory, err := suite.repository.Create(category)
		assert.NoError(t, err)
		assert.NotNil(t, createdCategory)

		if t.Failed() {
			t.Fatal("CreateCategory failed, stopping sequential execution")
		}
	})

	t.Run("GetAllCategoriesForUser", func(t *testing.T) {
		categories, err := suite.repository.GetAllForUser("1244ABC")

		assert.NoError(t, err)
		assert.NotNil(t, categories)
		assert.Len(t, categories, 2)
	})

	t.Run("UpdateCategory", func(t *testing.T) {
		category := &models.DBCategory{
			ID:     "2345SDSXAS",
			UserID: "1244ABC",
			Name:   "Do by today",
			Color:  "#ffff00",
		}
		suite.repository.Update(category)

		updated_category, err := suite.repository.GetById("2345SDSXAS")
		assert.NoError(t, err)
		assert.NotNil(t, updated_category)
		assert.Equal(t, "Do by today", updated_category.Name)
		assert.Equal(t, "#ffff00", updated_category.Color)
	})

	t.Run("DeleteCategory", func(t *testing.T) {
		err := suite.repository.Delete("2345SDSXAS")
		assert.NoError(t, err)

		_, err = suite.repository.GetById("2345SDSXAS")
		expectedErrorMessage := "Resource not found, sql: no rows in result set"
		assert.Contains(t, err.Error(), expectedErrorMessage)
	})
}

func TestCategoryRepoTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepoTestSuite))
}
