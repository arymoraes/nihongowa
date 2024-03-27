package config

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func OpenAIInit() {
	// Initialize the OpenAI client with the API key from environment variables
	OpenAIClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}
