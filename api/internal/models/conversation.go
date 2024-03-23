package models

import (
	"fmt"
	"nihongowa/internal/config"
	"time"

	"github.com/gocql/gocql"
)

type Conversation struct {
	ID       gocql.UUID `json:"id"`
	Messages []Message  `json:"messages"`
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
			message.Content = rawMessage["content"].(string)
			message.Translation = rawMessage["translation"].(string)
			message.WordByWordTranslation = rawMessage["wordbywordtranslation"].([]string)
			message.CreatedAt = rawMessage["createdat"].(time.Time)
			message.UpdatedAt = rawMessage["updatedat"].(time.Time)

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
