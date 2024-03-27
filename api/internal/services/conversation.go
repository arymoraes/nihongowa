package services

import (
	"fmt"
	"nihongowa/internal/models"

	"github.com/gocql/gocql"
)

func CreateConversation(conversation *models.Conversation) (gocql.UUID, error) {
	err := models.NewConversation(conversation)

	if err != nil {
		fmt.Println("Error creating conversation", err)
		return gocql.UUID{}, err
	}

	new_message := conversation.Messages[0]
	new_message.IsAI = true
	new_message.ConversationID = conversation.ID.String()

	models.NewMessage(&new_message)

	if err != nil {
		fmt.Println("Error creating conversation", err)
		return gocql.UUID{}, err
	}

	return conversation.ID, nil
}
