package models

import "time"

type User struct {
	ID        *string    `json:"id"`
	Username  *string    `json:"username,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
