package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kjj1998/task-management-system/internal/database"
	"github.com/kjj1998/task-management-system/internal/models"
	"github.com/kjj1998/task-management-system/internal/repository"
)

func main() {
	env := os.Getenv("ENV")

	// Check environment
	if env == "dev" {
		fmt.Println("running in development mode")

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// initialize repositories
	userRepo := repository.NewUserRepository(database.GetDb())

	// fmt.Println(database.GetDb())
	user := &models.DbUser{
		Email:        "email.com",
		PasswordHash: "12423jnj34",
		FirstName:    "John",
		LastName:     "Doe",
	}

	userID, err := userRepo.Create(user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added user: %v\n", userID)
}
