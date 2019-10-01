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

// Создать новый чат между пользователями
// Запрос: POST /chats/add {"name": "chat_1", "users": ["<USER_ID_1>", "<USER_ID_2>"]}
// Ответ: id созданного чата или HTTP-код ошибки. Количество пользователей не ограничено.
func AddChat(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.ChatAdd
	err := decoder.Decode(&data)

	if err != nil || data.Name == nil || *data.Name == "" || data.UsersIds == nil || len(*data.UsersIds) == 0 {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	db := databases.GetPostgresSession()

	for _, userId := range *data.UsersIds {
		row := db.QueryRow("SELECT id FROM \"User\" WHERE id = $1;", userId)
		var id string
		err = row.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(404)
				io.WriteString(w, "Нет пользователя с id "+userId)
				return
			} else {
				log.Println(err)
				http.Error(w, http.StatusText(500), 500)
				return
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

	for _, userId := range *data.UsersIds {
		_, err := db.Exec("INSERT INTO \"Chat_User\" VALUES($1, $2);", userId, chatId)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}

	w.WriteHeader(201)
	result, _ := json.Marshal(models.Chat{ID: &chatId})
	io.WriteString(w, string(result))
}
