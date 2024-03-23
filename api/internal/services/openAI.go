package services

import (
	"context"
	"encoding/json"
	"fmt"
	"nihongowa/internal/config"
	"nihongowa/internal/models"

	"github.com/sashabaranov/go-openai"
)

// type Message struct {
// 	Content               string    `json:"content"`
// 	Translation           string    `json:"translation"`
// 	WordByWordTranslation []string  `json:"wordByWordTranslation"`
// 	CreatedAt             time.Time `json:"createdAt"`
// 	UpdatedAt             time.Time `json:"updatedAt"`
// }

// func SendMessageToChatGPT(message string) string {

// 	var assitantName = "Nihongowa Assistant"
// 	var assistantDescription = "A Japanese language learning assistant"
// 	var assistantInstructions = "You will hold a conversation in Japanese. You will receive a message from the user, and you will answer in Japanese, like a normal conversation. You will only use hiragana and katana, kanjis are forbidden.\n" +
// 		"Your response will be in a JSON format, and you will not send anything other than the JSON with the response, so I can parse the JSON on my server from your response. This is the JSON format:\n" +
// 		"{\n" +
// 		"    \"content\": \"こんにちは\",\n" +
// 		"    \"translation\": \"Hello\",\n" +
// 		"    \"wordByWordTranslation\": {\"こんにちは\": \"Hello\"}\n" +
// 		"}"

// 	assistant, err := config.OpenAIClient.CreateAssistant(
// 		context.Background(),
// 		openai.AssistantRequest{
// 			Name:         &assitantName,
// 			Description:  &assistantDescription,
// 			Model:        openai.GPT4TurboPreview,
// 			Instructions: &assistantInstructions,
// 		},
// 	)

// 	config.OpenAIClient.CreateThread(
// 		context.Background(),
// 		openai.ThreadRequest{
// 			AssistantID: assistant.ID,
// 			Name:        "nihongowa",
// 		},
// 	)

// 	if err != nil {
// 		return "Error"
// 	}

// 	return resp.Choices[0].Message.Content
// }

func SendMessageToChatGPT(message string) (models.Message, error) {
	var instructions = "You will hold a conversation in Japanese. You will receive a message from the user, and you will answer in Japanese, like a normal conversation. You will only use hiragana and katakana, kanjis are forbidden. If you don't know how to avoid kanjis, translate it to romanji and then transform it into hiragana/katakana.\n" +
		"Your response will be in a JSON format, and you will not send anything other than the JSON with the response, so I can parse the JSON on my server from your response. This is the JSON format:\n" +
		"{\n" +
		"    \"content\": \"こんにちは\",\n" +
		"    \"translation\": \"Hello\",\n" +
		"}\n" +
		"Again, kanjis are forbidden!!! Do not use kanjis inside `content`"

	resp, err := config.OpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: instructions,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	if err != nil {
		fmt.Println("error:", err)
		return models.Message{}, err
	}

	var response_model models.Message

	response := resp.Choices[0].Message.Content

	println(response)

	unmarshal_err := json.Unmarshal([]byte(response), &response_model)

	if unmarshal_err != nil {
		fmt.Println("error:", unmarshal_err)
		return models.Message{}, unmarshal_err
	}

	return response_model, nil
}
