package models

type MessageAdd struct {
	ChatId   *string `json:"chat"`
	AuthorId *string `json:"author"`
	Text     *string `json:"text"`
}
