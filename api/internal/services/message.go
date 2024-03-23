package services

import (
	"fmt"
	"log"

	"nihongowa/internal/models"

	"github.com/gocql/gocql"
)

func GetMessagesFromConversation(conversationID string) ([]models.Message, error) {
	id, err := gocql.ParseUUID(conversationID)
	if err != nil {
		log.Printf("Error parsing conversation ID: %v", err)
		return nil, err
	}

	conversation := models.Conversation{ID: id}

	if err := conversation.GetMessages(); err != nil {
		log.Printf("Error getting messages from conversation: %v", err)
		return nil, err
	}

	fmt.Println(conversation.Messages)

	return conversation.Messages, nil
}

func PostMessageToConversation(conversationID string, message models.Message) error {
	id, err := gocql.ParseUUID(conversationID)
	if err != nil {
		log.Printf("Error parsing conversation ID: %v", err)
		return err
	}

	conversation := models.Conversation{ID: id}

	if err := conversation.AddMessage(message); err != nil {
		log.Printf("Error posting message to conversation: %v", err)
		return err
	}

	// Optionally, send the message content to ChatGPT or perform other logic as needed
	// aiReply := SendMessageToChatGPT(message.Content)
	// Consider how you'll use aiReply or integrate it into your application

	return nil
}
