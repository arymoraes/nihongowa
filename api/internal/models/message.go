package models

import (
	"time"
)

type Message struct {
	Content               string    `json:"content"`
	Translation           string    `json:"translation"`
	WordByWordTranslation []string  `json:"wordByWordTranslation"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}
