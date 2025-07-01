package repository_test

import (
	"context"
	"log"
	"testing"

	"github.com/kjj1998/task-management-system/internal/database"
	"github.com/kjj1998/task-management-system/internal/models"
	"github.com/kjj1998/task-management-system/internal/repository"
	"github.com/kjj1998/task-management-system/internal/repository/test/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepoTestSuite struct {
	suite.Suite
	mySQLContainer *testhelpers.MySQLContainer
	ctx            context.Context
	repository     repository.UserRepository
}

func (suite *UserRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	mySQLContainer, err := testhelpers.CreateMySQLContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.mySQLContainer = mySQLContainer
	host, _ := mySQLContainer.Container.Host(suite.ctx)
	port, _ := mySQLContainer.Container.MappedPort(suite.ctx, "3306")

	database.Connect2("testuser", "testpass", host, port.Port(), "taskapi")
	db := database.GetDb()
	userRepository := repository.NewUserRepository(db)
	suite.repository = userRepository
}

func (suite *UserRepoTestSuite) TearDownSuite() {
	if err := suite.mySQLContainer.Container.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating mysql container: %s", err)
	}
}

func (suite *UserRepoTestSuite) TestCreateUser() {
	t := suite.T()

	user_id, err := suite.repository.Create(&models.DBUser{
		Email: "jc@email.com", PasswordHash: "dsfdsf!@#!23", FirstName: "Julus", LastName: "Caesar",
	})
	assert.NoError(t, err)
	assert.NotNil(t, user_id)
}

func TestCustomerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
