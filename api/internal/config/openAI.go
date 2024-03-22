package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func OpenAIInit() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	OpenAIClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}
