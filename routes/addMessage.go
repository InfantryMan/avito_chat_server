package routes

import (
	"awesomeProject/databases"
	"awesomeProject/models"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func AddMessage(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.MessageAdd
	err := decoder.Decode(&data)

	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln(err)
	}

	row = db.QueryRow("SELECT id FROM \"User\" WHERE id = $1", data.AuthorId)
	var authorId string
	if err = row.Scan(&authorId); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			io.WriteString(w, "Пользователь с id = "+*data.AuthorId+" не найден.")
			return
		}
		log.Fatalln(err)
	}

	row = db.QueryRow("INSERT INTO \"Message\" VALUES(DEFAULT, $1, DEFAULT, $2, $3) RETURNING id;", data.Text, chatId, authorId)
	var messageId string
	if err = row.Scan(&messageId); err != nil {
		log.Fatalln(err)
	}

	w.WriteHeader(201)
	result, _ := json.Marshal(models.Message{ID: &messageId})
	io.WriteString(w, string(result))
}
