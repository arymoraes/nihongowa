package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func OpenAIInit() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	OpenAIClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}
