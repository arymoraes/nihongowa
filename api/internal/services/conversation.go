package services

import (
	"fmt"
	"nihongowa/internal/config"
	"nihongowa/internal/models"

	"github.com/gocql/gocql"
)

func CreateConversation(conversation *models.Conversation) (gocql.UUID, error) {
	id := gocql.TimeUUID()
	conversation.ID = id

	err := config.Session.Query("INSERT INTO conversations (id, assistantid, threadid, runid, messages) VALUES (?, ?, ?, ?, ?)", id.String(),
		conversation.AssistantID, conversation.ThreadID, conversation.RunID, []models.Message{}).Exec()

	new_message := conversation.Messages[0]
	new_message.IsAI = true

	conversation.AddMessage(new_message)

	if err != nil {
		fmt.Println("Error creating conversation", err)
		return gocql.UUID{}, err
	}

	return id, nil
}
