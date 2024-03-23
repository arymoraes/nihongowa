package services

import (
	"nihongowa/internal/config"
	"nihongowa/internal/models"

	"github.com/google/uuid"
)

func CreateConversation() (uuid.UUID, error) {
	id := uuid.New()

	err := config.Session.Query("INSERT INTO conversations (id, messages) VALUES (?, ?) IF NOT EXISTS", id.String(), []models.Message{}).Exec()

	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
