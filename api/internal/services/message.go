package services

import (
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

	return conversation.Messages, nil
}

func PostMessageToConversation(conversationID string, message models.Message) (models.Message, error) {
	id, err := gocql.ParseUUID(conversationID)
	if err != nil {
		log.Printf("Error parsing conversation ID: %v", err)
		return models.Message{}, err
	}

	conversation, err := models.GetConversationById(id)

	if err != nil {
		log.Printf("Error getting conversation by ID: %v", err)
		return models.Message{}, err
	}

	if err := conversation.AddMessage(message); err != nil {
		log.Printf("Error posting message to conversation: %v", err)
		return models.Message{}, err
	}

	// AI Reply
	aiReply, aiErr := SendMessageToChatGPT(message.Content, conversation)

	println(aiReply.Content)

	if aiErr != nil {
		log.Printf("Error sending message to ChatGPT: %v", aiErr)
		return models.Message{}, aiErr
	}

	conversation.AddMessage(aiReply)

	return aiReply, nil
}
