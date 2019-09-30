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

// Получение всех сообщений чата с id <CHAT_ID>
// Запрос: POST /messages/get {"chat": "<CHAT_ID>"}
// Ответ: список всех сообщений чата со всеми полями,
// отсортированный по времени создания сообщения (от раннего к позднему).
// Или HTTP-код ошибки.
func GetMessages(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var data models.MessagesGet
	err := decoder.Decode(&data)

	if err != nil || data.ChatId == nil || *data.ChatId == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	db := databases.GetPostgresSession()

	var chatId string
	err = db.QueryRow("SELECT id FROM \"Chat\" WHERE id = $1", data.ChatId).Scan(&chatId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			io.WriteString(w, "Чат с id "+*data.ChatId+" не найден.")
			return
		}
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rows, err := db.Query("SELECT id, author_id, text, created_at FROM \"Message\" M WHERE chat_id = $1 ORDER BY created_at;", chatId)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	messages := make([]*models.Message, 0)
	for rows.Next() {
		message := new(models.Message)
		message.ChatId = &chatId
		err := rows.Scan(&message.ID, &message.AuthorId, &message.Text, &message.CreatedAt)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	res, err := json.Marshal(messages)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(200)
	io.WriteString(w, string(res))
}
