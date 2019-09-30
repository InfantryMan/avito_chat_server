package routes

import (
	"awesomeProject/databases"
	"awesomeProject/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Добавить нового пользователя
// Запрос: POST /users/add {"username": "user_1"}
// Ответ: id созданного пользователя или HTTP-код ошибки.
func AddUser(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.UserAdd
	err := decoder.Decode(&data)

	if err != nil || data.Username == nil || *data.Username == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	dbSession := databases.GetPostgresSession()

	row := dbSession.QueryRow("SELECT id FROM \"User\" WHERE username=$1;", data.Username)
	var userId string
	if err = row.Scan(&userId); err == nil && userId != "" {
		w.WriteHeader(409)
		io.WriteString(w, "Пользователь с таким ником уже существует.")
		return
	}

	row = dbSession.QueryRow("INSERT INTO \"User\" VALUES(DEFAULT, $1) RETURNING id;", data.Username)
	if err = row.Scan(&userId); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	result, err := json.Marshal(models.User{ID: &userId})
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(201)
	io.WriteString(w, string(result))
}
