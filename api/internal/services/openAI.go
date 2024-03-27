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

var instructions = "You will hold a conversation in Japanese. You will receive a message from the user (will be the previous message on the thread), and you will answer in Japanese, like a normal conversation. You will try to be brief in your responses (1-2 sentences)" +
	"You will use hiragana and katakana, and avoid using kanjis\n" +
	"Your response will be in a JSON format, and you will not send anything other than the JSON with the response, so I can parse the JSON on my server from your response." +
	"The JSON will contain: `content`, `romanji` and `translation`. `content` will be the message in Japanese, and `translation` will be the translation of the message in English. You will also include `romanji` which will be the romanji of the message in Japanese." +
	"This is the JSON format:\n" +
	"{\n" +
	"    \"content\": \"こんにちは\",\n" +
	"    \"translation\": \"Hello\"\n" +
	"    \"romanji\": \"konnichiwa\"\n" +
	"    \"user_message_translated\": \"Hello\"\n" +
	"}\n" +
	"If the user sends a message in English or in romanji, you will respond in Japanese like normal. The conversation will be started with a user message of \"こんにちは\".\n" +
	"In which you will use the scenario I give you to create a conversation."

func CreateConversationScenario() (models.Conversation, error) {
	conversation := models.Conversation{}

	assistantId, err := createAssistant(&conversation)

	if err != nil {
		fmt.Println("Error creating assistant", err)
		return models.Conversation{}, err
	}

	conversation.AssistantID = assistantId

	run, err := createThread(conversation.AssistantID, "こんにちは")

	if err != nil {
		return models.Conversation{}, err
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

	message, err := retrieveAndProcessMessages(conversation.ThreadID)

	if err != nil {
		return models.Conversation{}, err
	}

	conversation.Messages = []models.Message{message}
	conversation.RunID = run.ID

	_, err = CreateConversation(&conversation)

	if err != nil {
		fmt.Println("Error creating conversation", err)
		return models.Conversation{}, err
	}

	return conversation, nil
}

func SendMessageToChatGPT(message string, conversation *models.Conversation) (models.Message, error) {
	// Send a message
	_, err := config.OpenAIClient.CreateMessage(context.Background(), conversation.ThreadID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	if err != nil {
		fmt.Println("Error creating message", err)
		return models.Message{}, err
	}

	scenario := conversation.Scenario

	run, err := config.OpenAIClient.CreateRun(context.Background(), conversation.ThreadID, openai.RunRequest{
		AssistantID:  conversation.AssistantID,
		Model:        openai.GPT3Dot5Turbo,
		Instructions: instructions + "Scenario: " + scenario + "\n",
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

	return retrieveAndProcessMessages(conversation.ThreadID)
}

type Scenario struct {
	Scenario string `json:"scenario"`
}

func createAssistant(c *models.Conversation) (string, error) {
	// Generate Assistant Name
	names := []string{}

	fmt.Println("Base path", config.BasePath)

	filePath := fmt.Sprintf("%snames.json", config.BasePath)
	namesFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println("Error opening names file", err)
		return "", err
	}

	namesByteValue, _ := ioutil.ReadAll(namesFile)
	json.Unmarshal(namesByteValue, &names)
	names_n := rand.Int() % len(names)

	assistantName := names[names_n]

	defer namesFile.Close()

	var name = assistantName
	var description = "A Japanese language learning assistant"

	// Generate Scenario
	var scenarios []Scenario

	filePath = fmt.Sprintf("%sscenarios.json", config.BasePath)
	scenariosFile, err := os.Open(filePath)

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
	c.AssistantName = assistantName

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
