package models

type ChatAdd struct {
	Name     *string  `json:"name"`
	UsersIds []string `json:"users"`
}
