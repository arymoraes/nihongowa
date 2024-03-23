package services

import (
	"context"
	"encoding/json"
	"fmt"
	"nihongowa/internal/config"
	"nihongowa/internal/models"
	"time"

	"github.com/sashabaranov/go-openai"
)

func SendMessageToChatGPT(message string, conversation *models.Conversation) (models.Message, error) {
	if conversation.AssistantID == "" {
		assistantId, err := createAssistant()

		if err != nil {
			fmt.Println("Error creating assistant", err)
			return models.Message{}, err
		}

		conversation.AssistantID = assistantId
	}

	if conversation.ThreadID == "" {
		// In case the thread doesn't exist, we create a new one
		run, err := createThread(conversation.AssistantID, message)

		if err != nil {
			return models.Message{}, err
		}

		conversation.ThreadID = run.ThreadID

		for run.Status == openai.RunStatusQueued || run.Status == openai.RunStatusInProgress {
			response, err := config.OpenAIClient.RetrieveRun(context.Background(), run.ThreadID, run.ID)

			if err != nil {
				fmt.Println("Error retrieving run", err)
			}

			run = response
			time.Sleep(1000 * time.Millisecond)
		}
	} else {
		// If the thread exists, we spawn a run
		run, err := config.OpenAIClient.CreateRun(context.Background(), conversation.ThreadID, openai.RunRequest{
			AssistantID: conversation.AssistantID,
			Model:       openai.GPT3Dot5Turbo,
		})

		if err != nil {
			fmt.Println("Error creating run", err)
			return models.Message{}, err
		}

		for run.Status == openai.RunStatusQueued || run.Status == openai.RunStatusInProgress {
			response, err := config.OpenAIClient.RetrieveRun(context.Background(), run.ThreadID, run.ID)

			if err != nil {
				fmt.Println("Error retrieving run", err)
			}

			run = response
			time.Sleep(1000 * time.Millisecond)
		}
	}

	return retrieveAndProcessMessages(conversation.ThreadID)
}

func createAssistant() (string, error) {
	var name = "Nihongowa Assistant"
	var description = "A Japanese language learning assistant"
	var instructions = "You will hold a conversation in Japanese. You will receive a message from the user, and you will answer in Japanese, like a normal conversation. You will only use hiragana and katakana, kanjis are forbidden. If you don't know how to avoid kanjis, translate it to romanji and then transform it into hiragana/katakana.\n" +
		"Your response will be in a JSON format, and you will not send anything other than the JSON with the response, so I can parse the JSON on my server from your response. This is the JSON format:\n" +
		"{\n" +
		"    \"content\": \"こんにちは\",\n" +
		"    \"translation\": \"Hello\"\n" +
		"}\n" +
		"Again, kanjis are forbidden!!! Do not use kanjis inside `content`"

	assistant, err := config.OpenAIClient.CreateAssistant(
		context.Background(),
		openai.AssistantRequest{
			Name:         &name,
			Description:  &description,
			Model:        openai.GPT3Dot5Turbo,
			Instructions: &instructions,
		},
	)

	if err != nil {
		fmt.Println("Error creating assistant", err)
		return "", err
	}

	return assistant.ID, nil
}

func createThread(assistantId string, message string) (openai.Run, error) {
	run, err := config.OpenAIClient.CreateThreadAndRun(context.Background(), openai.CreateThreadAndRunRequest{
		RunRequest: openai.RunRequest{AssistantID: assistantId},
		Thread: openai.ThreadRequest{
			Messages: []openai.ThreadMessage{{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			}},
		},
	})

	if err != nil {
		fmt.Println("Error creating thread", err)
		return openai.Run{}, err
	}

	return run, nil
}

// func createRunInThread(conversation *models.Conversation) error {
// 	run, err := config.OpenAIClient.CreateRun(context.Background(), conversation.ThreadID, openai.RunRequest{
// 		AssistantID: conversation.AssistantID,
// 		Model:       openai.GPT3Dot5Turbo,
// 	})
// 	if err != nil {
// 		fmt.Println("Error creating run", err)
// 		return err
// 	}

// 	return waitForRunCompletion(&run)
// }

func retrieveAndProcessMessages(threadID string) (models.Message, error) {
	numMessages := 1
	messages, err := config.OpenAIClient.ListMessage(context.Background(), threadID, &numMessages, nil, nil, nil)
	if err != nil {
		fmt.Println("Error listing messages", err)
		return models.Message{}, err
	}

	if len(messages.Messages) == 0 {
		return models.Message{}, fmt.Errorf("no messages found in thread")
	}

	var responseModel models.Message
	response := messages.Messages[0].Content[0].Text.Value
	unmarshalErr := json.Unmarshal([]byte(response), &responseModel)
	if unmarshalErr != nil {
		fmt.Println("error:", unmarshalErr)
		return models.Message{}, unmarshalErr
	}

	return responseModel, nil
}

// func waitForRunCompletion(run *openai.Run) error {
// 	for run.Status == openai.RunStatusQueued || run.Status == openai.RunStatusInProgress {
// 		response, err := config.OpenAIClient.RetrieveRun(context.Background(), run.ThreadID, run.ID)
// 		if err != nil {
// 			fmt.Println("Error retrieving run", err)
// 			return err
// 		}
// 		run = &response
// 		time.Sleep(1000 * time.Millisecond)
// 	}
// 	return nil
// }
