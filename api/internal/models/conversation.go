package models

import (
	"fmt"
	"nihongowa/internal/config"
	"time"

	"github.com/gocql/gocql"
)

type Conversation struct {
	ID            gocql.UUID `json:"id"`
	Messages      []Message  `json:"messages"`
	AssistantName string     `json:"assistant_name"`
	RunID         string     `json:"run_id"`
	ThreadID      string     `json:"thread_id"`
	AssistantID   string     `json:"assistant_id"`
	Scenario      string     `json:"scenario"`
	LastMessageAt time.Time  `json:"last_message_at"`
}

func NewConversation(conversation *Conversation) error {
	id := gocql.TimeUUID()
	conversation.ID = id

	err := config.Session.Query("INSERT INTO conversations (id, assistant_id, thread_id, run_id, scenario, assistant_name) VALUES (?, ?, ?, ?, ?, ?)", id.String(),
		conversation.AssistantID, conversation.ThreadID, conversation.RunID, conversation.Scenario, conversation.AssistantName).Exec()

	if err != nil {
		fmt.Println("Error creating conversation", err)
		return err
	}

	return nil
}

func GetConversationById(conversationID gocql.UUID) (*Conversation, error) {
	conversation := &Conversation{ID: conversationID}

	err := config.Session.Query(`SELECT id, assistant_name, thread_id, assistant_id, scenario, run_id FROM conversations WHERE id = ? LIMIT 1`,
		conversation.ID).Consistency(gocql.One).Scan(&conversation.ID, &conversation.AssistantName,
		&conversation.ThreadID, &conversation.AssistantID, &conversation.Scenario, &conversation.RunID)

	if err != nil {
		return nil, err
	}

	return conversation, nil
}

func UpdateConversation(conversation *Conversation) error {
	query := "UPDATE conversations SET thread_id = ?, assistant_id = ?, scenario = ? WHERE id = ?"

	if err := config.Session.Query(query, conversation.ThreadID, conversation.AssistantID, conversation.Scenario, conversation.ID).Exec(); err != nil {
		return err
	}

	return nil
}

func GetLastConversations() ([]Conversation, error) {
	iter := config.Session.Query("SELECT id, assistant_name FROM conversations LIMIT 10").Iter()

	conversations := []Conversation{}

	m := map[string]interface{}{}
	for iter.MapScan(m) {
		conversationID, ok := m["id"].(gocql.UUID)
		if !ok {
			continue
		}

		conversation, err := GetConversationById(conversationID)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, *conversation)

		m = map[string]interface{}{}
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return conversations, nil
}
