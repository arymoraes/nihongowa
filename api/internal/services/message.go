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

	messages, err := models.GetMessagesByConversationID(id.String())

	if err != nil {
		log.Printf("Error getting messages by conversation ID: %v", err)
		return nil, err
	}

	return messages, nil
}

func PostMessageToConversation(conversationID string, message models.Message) (models.Message, error) {
	id, err := gocql.ParseUUID(conversationID)
	if err != nil {
		log.Printf("Error parsing conversation ID: %v", err)
		return models.Message{}, err
	}

	conversation, err := models.GetConversationById(id)

	if err != nil {
		log.Printf("Error getting conversation by ID: %v. ID: %v", err, id)
		return models.Message{}, err
	}

	message.ConversationID = conversationID

	if err := models.NewMessage(&message); err != nil {
		log.Printf("Error posting message to conversation: %v", err)
		return models.Message{}, err
	}

	// AI Reply
	aiReply, aiErr := SendMessageToChatGPT(message.Content, conversation)

	if aiErr != nil {
		log.Printf("Error sending message to ChatGPT: %v", aiErr)
		return models.Message{}, aiErr
	}

	aiReply.ConversationID = conversationID
	aiReply.IsAI = true

	models.NewMessage(&aiReply)

	return aiReply, nil
}

func GetLastConversations() ([]models.Conversation, error) {
	conversations, err := models.GetLastConversations()

	if err != nil {
		log.Printf("Error getting last conversations: %v", err)
		return nil, err
	}

	return conversations, nil
}
