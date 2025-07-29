package user_test

import (
	"context"
	"log"
	"testing"

	"github.com/kjj1998/task-management-system/internal/database"
	"github.com/kjj1998/task-management-system/internal/errors"
	"github.com/kjj1998/task-management-system/internal/logger"
	"github.com/kjj1998/task-management-system/internal/models"
	"github.com/kjj1998/task-management-system/internal/repository/testutils"
	"github.com/kjj1998/task-management-system/internal/repository/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepoTestSuite struct {
	suite.Suite
	mySQLContainer *testutils.MySQLContainer
	ctx            context.Context
	repository     user.UserRepository
}

func (suite *UserRepoTestSuite) SetupSuite() {
	logger := logger.NewLogger("test")
	suite.ctx = context.Background()

	mySQLContainer, err := testutils.CreateMySQLContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.mySQLContainer = mySQLContainer
	host, _ := mySQLContainer.Container.Host(suite.ctx)
	port, _ := mySQLContainer.Container.MappedPort(suite.ctx, "3306")

	err = database.Connect("testuser", "testpass", host, port.Port(), "taskapi", logger)
	suite.Require().NoError(err, "Failed to connect to test database")
	db := database.GetDb()
	dbErrorHandler := errors.NewDatabaseErrorHandler()
	userRepository := user.NewUserRepository(db, dbErrorHandler, logger)
	suite.repository = userRepository
}

func (suite *UserRepoTestSuite) TearDownSuite() {
	if err := suite.mySQLContainer.Container.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating mysql container: %s", err)
	}
}

func (suite *UserRepoTestSuite) TestUserLifecycle() {
	t := suite.T()

	t.Run("CreateUser", func(t *testing.T) {
		user := &models.DBUser{
			Email:        "jc@email.com",
			PasswordHash: "dsfdsf!@#!23",
			FirstName:    "Julius",
			LastName:     "Caesar",
		}

		created_user, err := suite.repository.Create(user)
		assert.NoError(t, err)
		assert.NotNil(t, created_user)

		if t.Failed() {
			t.Fatal("CreateUser failed, stopping sequential execution")
		}
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		user, err := suite.repository.GetByEmail("john@email.com")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john@email.com", user.Email)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		user, err := suite.repository.GetByEmail("john@email.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)

		err = suite.repository.Update(&models.DBUser{Email: "johnathan@email.com", FirstName: "Johnathan", LastName: "Doe", ID: user.ID})
		assert.NoError(t, err)

		user, err = suite.repository.GetById(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "Johnathan", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "johnathan@email.com", user.Email)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		err := suite.repository.Delete("1244ABC")
		assert.NoError(t, err)

		_, err = suite.repository.GetById("1244ABC")
		expectedErrorMessage := "Resource not found, sql: no rows in result set"
		assert.Contains(t, err.Error(), expectedErrorMessage)
	})
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
