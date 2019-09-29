package models

import "time"

type Chat struct {
	ID        *string    `json:"id"`
	Name      *string    `json:"name,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
