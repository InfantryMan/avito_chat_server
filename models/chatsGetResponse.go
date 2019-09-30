package models

import "time"

type ChatsGetResponse struct {
	ID        *string    `json:"id"`
	Name      *string    `json:"name"`
	Users     *[]string  `json:"users"`
	CreatedAt *time.Time `json:"created_at"`
}
