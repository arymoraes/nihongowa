package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"nihongowa/internal/config"
	"nihongowa/internal/models"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

var instructions = "You will hold a conversation in Japanese. You will receive a message from the user (will be the previous message on the thread), and you will answer in Japanese, like a normal conversation." +
	"You will only use hiragana and katakana, kanjis are forbidden.\n" +
	"Your response will be in a JSON format, and you will not send anything other than the JSON with the response, so I can parse the JSON on my server from your response." +
	"The JSON will contain: `content` and `translation`. `content` will be the message in Japanese, and `translation` will be the translation of the message in English. You will also include `romanji` which will be the romanji of the message in Japanese." +
	"This is the JSON format:\n" +
	"{\n" +
	"    \"content\": \"こんにちは\",\n" +
	"    \"translation\": \"Hello\"\n" +
	"    \"romanji\": \"konnichiwa\"\n" +
	"    \"user_message_translated\": \"Hello\"\n" +
	"}\n" +
	"If the user sends a message in English or in romanji, you will respond in Japanese like normal"

func SendMessageToChatGPT(message string, conversation *models.Conversation) (models.Message, error) {
	if conversation.AssistantID == "" {
		assistantId, err := createAssistant(conversation)

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

		// We save the new ThreadID and AssistantID in the conversation
		models.UpdateConversation(conversation)

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
			AssistantID:  conversation.AssistantID,
			Model:        openai.GPT3Dot5Turbo,
			Instructions: instructions,
		})

		if err != nil {
			fmt.Println("Error creating run", err)
			return models.Message{}, err
		}

		// Send a message
		_, err = config.OpenAIClient.CreateMessage(context.Background(), conversation.ThreadID, openai.MessageRequest{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		})

		if err != nil {
			fmt.Println("Error creating message", err)
			return models.Message{}, err
		}

		for run.Status == openai.RunStatusQueued || run.Status == openai.RunStatusInProgress {
			response, err := config.OpenAIClient.RetrieveRun(context.Background(), run.ThreadID, run.ID)

			fmt.Println("run", run)

			if err != nil {
				fmt.Println("Error retrieving run", err)
			}

			run = response
			time.Sleep(1000 * time.Millisecond)
		}
	}

	return retrieveAndProcessMessages(conversation.ThreadID)
}

type Scenario struct {
	Scenario string `json:"scenario"`
}

func createAssistant(c *models.Conversation) (string, error) {
	var name = "Nihongowa Assistant"
	var description = "A Japanese language learning assistant"

	var scenarios []Scenario

	scenariosFile, err := os.Open("scenarios.json")

	if err != nil {
		fmt.Println("Error opening scenarios file", err)
		return "", err
	}

	byteValue, _ := ioutil.ReadAll(scenariosFile)
	json.Unmarshal(byteValue, &scenarios)
	n := rand.Int() % len(scenarios)

	scenario := scenarios[n].Scenario
	instructions += "Scenario: " + scenario + "\n"
	c.Scenario = scenario

	defer scenariosFile.Close()

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

	fmt.Println("messages", messages)

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
