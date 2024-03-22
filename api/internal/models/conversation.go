package models

import "github.com/gocql/gocql"

type Conversation struct {
	ID       gocql.UUID `json:"id"`
	Messages []Message  `json:"messages"`
}
