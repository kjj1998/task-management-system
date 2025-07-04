package user_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/kjj1998/task-management-system/internal/database"
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
	userRepository := user.NewUserRepository(db)
	suite.repository = userRepository
}

func (suite *UserRepoTestSuite) TearDownSuite() {
	if err := suite.mySQLContainer.Container.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating mysql container: %s", err)
	}
}

func (suite *UserRepoTestSuite) TestUserLifecycle() {
	t := suite.T()

	user_id, err := suite.repository.Create(&models.DBUser{
		Email: "jc@email.com", PasswordHash: "dsfdsf!@#!23", FirstName: "Julius", LastName: "Caesar",
	})
	assert.NoError(t, err)
	assert.NotNil(t, user_id)
}

func (suite *UserRepoTestSuite) TestGetUserByEmail() {
	t := suite.T()

	user, err := suite.repository.GetByEmail("john@email.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "john@email.com", user.Email)
}

func (suite *UserRepoTestSuite) TestUpdateUser() {
	t := suite.T()

	user, err := suite.repository.GetByEmail("john@email.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)

	suite.repository.Update(&models.DBUser{Email: "johnathan@email.com", FirstName: "Johnathan", LastName: "Doe", ID: user.ID})

	user, err = suite.repository.GetById(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Johnathan", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "johnathan@email.com", user.Email)
}

func (suite *UserRepoTestSuite) TestDeleteUser() {
	t := suite.T()

	user_id, err := suite.repository.Create(&models.DBUser{
		Email: "jason@email.com", PasswordHash: "dsfdsf!@#!23", FirstName: "Jason", LastName: "Bourne",
	})
	assert.NoError(t, err)
	assert.NotNil(t, user_id)

	user, err := suite.repository.GetByEmail("jason@email.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)

	err = suite.repository.Delete(user.ID)
	assert.NoError(t, err)

	_, err = suite.repository.GetById(user.ID)
	expectedErrorMessage := fmt.Sprintf("usersById %s: no such user", user_id)
	assert.Contains(t, err.Error(), expectedErrorMessage)
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
