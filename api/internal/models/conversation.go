package models

import (
	"fmt"
	"nihongowa/internal/config"
	"time"

	"github.com/gocql/gocql"
)

type Conversation struct {
	ID          gocql.UUID `json:"id"`
	Messages    []Message  `json:"messages"`
	ThreadID    string     `json:"thread_id"`
	AssistantID string     `json:"assistant_id"`
	Scenario    string     `json:"scenario"`
}

func GetConversationById(conversationID gocql.UUID) (*Conversation, error) {
	conversation := &Conversation{ID: conversationID}

	// err := config.Session.Query("SELECT COUNT(*) FROM conversations WHERE id = ?", conversationID).Consistency(gocql.One).Scan(&conversation)
	err := config.Session.Query(`SELECT id FROM conversations WHERE id = ? LIMIT 1`,
		conversation.ID).Consistency(gocql.One).Scan(&conversation.ID)

	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func UpdateConversation(conversation *Conversation) error {
	query := "UPDATE conversations SET threadid = ?, assistantid = ?, scenario = ? WHERE id = ?"

	if err := config.Session.Query(query, conversation.ThreadID, conversation.AssistantID, conversation.Scenario, conversation.ID).Exec(); err != nil {
		return err
	}

	return nil
}

func NewConversation(conversationID gocql.UUID) (*Conversation, error) {
	conversation := &Conversation{ID: conversationID, Messages: []Message{}}

	applied, err := config.Session.Query("INSERT INTO conversations (id, messages) VALUES (?, ?) IF NOT EXISTS",
		conversationID, conversation.Messages).Consistency(gocql.One).ScanCAS(&conversation.ID, &conversation.Messages)

	if err != nil {
		return nil, err
	}

	if !applied {
		return nil, fmt.Errorf("conversation already exists")
	}

	return conversation, nil
}

func (c *Conversation) AddMessage(message Message) error {
	messageMap := map[string]interface{}{
		"content":               message.Content,
		"translation":           message.Translation,
		"romanji":               message.Romanji,
		"isai":                  message.IsAI,
		"usermessagetranslated": message.UserMessageTranslated,
		"wordbywordtranslation": []string{},
		"createdat":             message.CreatedAt,
		"updatedat":             message.UpdatedAt,
	}

	query := "UPDATE conversations SET messages = messages + ? WHERE id = ?"
	if err := config.Session.Query(query, []map[string]interface{}{messageMap}, c.ID).Exec(); err != nil {
		return err
	}

	c.Messages = append(c.Messages, message)

	return nil
}

func (c *Conversation) GetMessages() error {
	var count int64
	err := config.Session.Query("SELECT COUNT(*) FROM conversations WHERE id = ?", c.ID).Consistency(gocql.One).Scan(&count)

	if err != nil {
		return err
	}

	if count == 0 {
		c.Messages = nil
		return nil
	}

	iter := config.Session.Query("SELECT messages FROM conversations WHERE id = ?", c.ID).Consistency(gocql.One).Iter()

	messages := []Message{}

	m := map[string]interface{}{}
	for iter.MapScan(m) {
		rawMessages, ok := m["messages"].([]map[string]interface{})
		if !ok {
			continue
		}

		for _, rawMessage := range rawMessages {
			var message Message

			// For strings that might be null
			if content, ok := rawMessage["content"].(string); ok {
				message.Content = content
			}
			if translation, ok := rawMessage["translation"].(string); ok {
				message.Translation = translation
			}
			if romanji, ok := rawMessage["romanji"].(string); ok {
				message.Romanji = romanji
			}
			if userMessageTranslated, ok := rawMessage["usermessagetranslated"].(string); ok {
				message.UserMessageTranslated = userMessageTranslated
			}

			// For slices of strings
			if wbt, ok := rawMessage["wordbywordtranslation"].([]interface{}); ok {
				for _, v := range wbt {
					if str, ok := v.(string); ok {
						message.WordByWordTranslation = append(message.WordByWordTranslation, str)
					}
				}
			}

			// For boolean fields
			if isAI, ok := rawMessage["isai"].(bool); ok {
				message.IsAI = isAI
			}

			// For time.Time fields assuming they are provided as string and might be null
			// If they're in a different format or you use *time.Time, this will need adjustment
			if createdAtStr, ok := rawMessage["createdat"].(string); ok && createdAtStr != "" {
				createdAt, err := time.Parse(time.RFC3339, createdAtStr) // Adjust the format as necessary
				if err == nil {
					message.CreatedAt = createdAt
				}
			}
			if updatedAtStr, ok := rawMessage["updatedat"].(string); ok && updatedAtStr != "" {
				updatedAt, err := time.Parse(time.RFC3339, updatedAtStr) // Adjust the format as necessary
				if err == nil {
					message.UpdatedAt = updatedAt
				}
			}

			messages = append(messages, message)
		}

		m = map[string]interface{}{}
	}

	if err := iter.Close(); err != nil {
		return err
	}

	c.Messages = messages

	return nil
}
