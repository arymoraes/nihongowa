package services

import (
	"fmt"
	"time"

	"nihongowa/internal/config"
	"nihongowa/internal/models"

	"github.com/gocql/gocql"
)

func GetMessagesFromConversation(conversationID string) []models.Message {
	conversation := models.Conversation{}

	iter := config.Session.Query("SELECT messages FROM conversations WHERE id = ?", conversationID).Consistency(gocql.One).Iter()

	var messages []models.Message

	// Use MapScan to iterate over the rows and unmarshal each message into a struct
	m := map[string]interface{}{}
	for iter.MapScan(m) {
		// Assuming 'messages' is the column name of the list of messages in Cassandra
		rawMessages, ok := m["messages"].([]map[string]interface{})
		if !ok {
			// handle the error: the type assertion didn't hold
			continue
		}

		for _, rawMessage := range rawMessages {
			var message models.Message
			// Convert the map to your Message struct
			// You may need to manually assign the fields depending on your exact schema
			message.Content = rawMessage["content"].(string)
			message.Translation = rawMessage["translation"].(string)
			message.WordByWordTranslation = rawMessage["wordbywordtranslation"].([]string)
			message.CreatedAt = rawMessage["createdat"].(time.Time)
			message.UpdatedAt = rawMessage["updatedat"].(time.Time)

			// Append the message to your slice
			messages = append(messages, message)
		}

		// Clear the map for the next iteration
		m = map[string]interface{}{}
	}

	// Check for any errors that occurred during the iteration
	if err := iter.Close(); err != nil {
		// handle the error
	}

	// Assign the messages to the conversation
	conversation.Messages = messages

	// Now, conversation.Messages should hold the unmarshalled messages from Cassandra
	fmt.Println(conversation.Messages)

	return conversation.Messages
}

func PostMessageToConversation(conversationID string, message models.Message) {
	// Create conversation if it doesn't exist
	config.Session.Query("INSERT INTO conversations (id, messages) VALUES (?, ?) IF NOT EXISTS", conversationID, []models.Message{}).Exec()

	messageMap := map[string]interface{}{
		"content":               message.Content,
		"translation":           message.Translation,
		"wordbywordtranslation": message.WordByWordTranslation,
		"createdat":             message.CreatedAt,
		"updatedat":             message.UpdatedAt,
	}

	query := "UPDATE conversations SET messages = messages + ? WHERE id = ?"
	err := config.Session.Query(query, []map[string]interface{}{messageMap}, conversationID).Exec()

	if err != nil {
		fmt.Println(err)
	}
}
