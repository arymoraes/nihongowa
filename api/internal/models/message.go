package models

import (
	"log"
	"nihongowa/internal/config"
	"time"

	"github.com/gocql/gocql"
)

type Message struct {
	ID                    string    `json:"id"`
	Content               string    `json:"content"`
	Translation           string    `json:"translation"`
	Romanji               string    `json:"romanji"`
	ConversationID        string    `json:"conversation_id"`
	UserMessageTranslated string    `json:"user_message_translated"`
	IsAI                  bool      `json:"is_ai"`
	WordByWordTranslation []string  `json:"word_by_word_translation"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func GetMessagesByConversationID(conversationId string) ([]Message, error) {
	scanner := config.Session.Query("SELECT * FROM messages WHERE conversation_id = ?", conversationId).Iter().Scanner()

	messages := []Message{}

	for scanner.Next() {
		var message Message

		err := scanner.Scan(&message.ConversationID, &message.ID, &message.Content, &message.CreatedAt, &message.IsAI, &message.Romanji,
			&message.Translation, &message.UpdatedAt, &message.UserMessageTranslated, &message.WordByWordTranslation)

		if err != nil {
			log.Fatal(err)
		}

		messages = append(messages, message)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return messages, nil
}

func NewMessage(message *Message) error {
	id := gocql.TimeUUID()
	created_at := time.Now()

	err := config.Session.Query("INSERT INTO messages (conversation_id, id, content, created_at, is_ai, romanji, translation, user_message_translated) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		message.ConversationID, id, message.Content, created_at, message.IsAI, message.Romanji, message.Translation, message.UserMessageTranslated).Exec()

	message.ID = id.String()

	if err != nil {
		log.Println("Error creating message", err)
		return err
	}

	return nil
}
