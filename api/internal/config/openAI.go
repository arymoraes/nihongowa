package config

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func openAIInit() {
	OpenAIClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}
