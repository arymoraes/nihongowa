package config

import (
	"os"

	"github.com/labstack/gommon/log"
	"github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func openAIInit() {
	log.Info("Initializing OpenAI client")
	OpenAIClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}
