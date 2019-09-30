package models

import "time"

type Message struct {
	ID        *string    `json:"id"`
	ChatId    *string    `json:"chat,omitempty"`
	AuthorId  *string    `json:"author,omitempty"`
	Text      *string    `json:"text,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
