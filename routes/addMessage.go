package routes

import (
	"database/sql"
	"encoding/json"
	"github.com/InfantryMan/avito_chat_server/databases"
	"github.com/InfantryMan/avito_chat_server/models"
	"io"
	"log"
	"net/http"
)

// Отправить сообщение в чат от лица пользователя
// Запрос: POST /messages/add {"chat": "<CHAT_ID>", "author": "<USER_ID>", "text": "hi"}
// Ответ: id созданного сообщения или HTTP-код ошибки.
func AddMessage(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.MessageAdd
	err := decoder.Decode(&data)

	if err != nil || data.ChatId == nil || *data.ChatId == "" ||
		data.AuthorId == nil || *data.AuthorId == "" ||
		data.Text == nil || *data.Text == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	db := databases.GetPostgresSession()

	row := db.QueryRow("SELECT id FROM \"Chat\" WHERE id = $1", data.ChatId)
	var chatId string
	if err = row.Scan(&chatId); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			io.WriteString(w, "Чата с id = "+*data.ChatId+" не существует.")
			return
		}
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	row = db.QueryRow("SELECT id FROM \"User\" WHERE id = $1", data.AuthorId)
	var authorId string
	if err = row.Scan(&authorId); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			io.WriteString(w, "Пользователь с id = "+*data.AuthorId+" не найден.")
			return
		}
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	row = db.QueryRow("SELECT user_id FROM \"Chat_User\" WHERE user_id = $1 AND chat_id = $2;", authorId, chatId)
	if err = row.Scan(&authorId); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			io.WriteString(w, "Пользователь с id = "+*data.AuthorId+" не состоит в таком чате.")
			return
		}
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	row = db.QueryRow("INSERT INTO \"Message\" VALUES(DEFAULT, $1, DEFAULT, $2, $3) RETURNING id;", data.Text, chatId, authorId)
	var messageId string
	if err = row.Scan(&messageId); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	result, err := json.Marshal(models.Message{ID: &messageId})
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(201)
	io.WriteString(w, string(result))
}
