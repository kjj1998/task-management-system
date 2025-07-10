package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kjj1998/task-management-system/internal/server"
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

	server := server.NewTaskManagementSystemServer()
	log.Fatal(http.ListenAndServe(":8080", server))
}
