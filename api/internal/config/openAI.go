package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func OpenAIInit() {
	// Check if the application is running in a Docker environment
	if os.Getenv("ENVIRONMENT") != "Docker" {
		// Attempt to load .env file if not running in Docker
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}

	// Initialize the OpenAI client with the API key from environment variables
	OpenAIClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}
