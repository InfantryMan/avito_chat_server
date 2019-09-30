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

func AddChat(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.ChatAdd
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	db := databases.GetPostgresSession()

	for _, userId := range data.UsersIds {
		row := db.QueryRow("SELECT id FROM \"User\" WHERE id = $1;", userId)
		var id string
		err = row.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(404)
				io.WriteString(w, "Нет пользователя с id "+userId)
				return
			} else {
				log.Fatalln(err)
			}
		}

	}

	row := db.QueryRow("INSERT INTO \"Chat\" VALUES(DEFAULT, $1) RETURNING id;", data.Name)
	var chatId string
	err = row.Scan(&chatId)
	if err != nil {
		w.WriteHeader(409)
		io.WriteString(w, "Чат с таким названием уже существует")
		return
	}

	for _, userId := range data.UsersIds {
		db.Exec("INSERT INTO \"Chat_User\" VALUES($1, $2);", userId, chatId)
	}

	w.WriteHeader(201)
	result, _ := json.Marshal(models.Chat{ID: &chatId})
	io.WriteString(w, string(result))
}
