package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
}
