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

// Получить список чатов конкретного пользователя
// Запрос: POST /chats/get {"user": "<USER_ID>"}
// Ответ: список всех чатов со всеми полями, отсортированный по времени создания последнего сообщения в чате
// (от позднего к раннему). Или HTTP-код ошибки.
func GetChats(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var data models.ChatsGet
	err := decoder.Decode(&data)

	if err != nil || data.UserId == nil || *data.UserId == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	db := databases.GetPostgresSession()

	row := db.QueryRow("SELECT id FROM \"User\" WHERE id = $1;", data.UserId)
	var userId string
	if err = row.Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			io.WriteString(w, "Нет пользователя с id "+*data.UserId)
			return
		}
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	query := "SELECT M.chat_id FROM \"Message\" M WHERE author_id=$1 GROUP BY M.chat_id ORDER BY max(M.created_at) DESC;"
	rows, err := db.Query(query, userId)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	chatIds := []string{}
	var chatId string
	for rows.Next() {
		if err := rows.Scan(&chatId); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		chatIds = append(chatIds, chatId)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	resultChats := []models.ChatsGetResponse{}
	var resultChat models.ChatsGetResponse

	for _, chatId = range chatIds {
		resultChat = models.ChatsGetResponse{ID: &chatId}
		row = db.QueryRow("SELECT name, created_at FROM \"Chat\" WHERE id=$1;", chatId)
		if err = row.Scan(&resultChat.Name, &resultChat.CreatedAt); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		rows, err := db.Query("SELECT U.username FROM \"Chat_User\" CU JOIN \"User\" U ON (CU.user_id=U.id) WHERE CU.chat_id = $1;", chatId)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		defer rows.Close()

		var usernames []string
		var username string
		for rows.Next() {
			if err := rows.Scan(&username); err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), 500)
				return
			}
			usernames = append(usernames, username)
		}

		if err = rows.Err(); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		resultChat.Users = &usernames
		resultChats = append(resultChats, resultChat)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	res, err := json.Marshal(resultChats)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(200)
	io.WriteString(w, string(res))
}
