package routes

import (
	"awesomeProject/databases"
	"awesomeProject/models"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func AddUser(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.UserAdd
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	dbSession := databases.GetPostgresSession()

	row := dbSession.QueryRow("INSERT INTO \"User\" VALUES(DEFAULT, $1) RETURNING id", data.Username)
	var id string
	if err = row.Scan(&id); err == nil {
		w.WriteHeader(201)
		result, _ := json.Marshal(models.User{ID: &id})
		io.WriteString(w, string(result))
	} else {
		w.WriteHeader(409)
		io.WriteString(w, "Пользователь с таким ником уже существует.")
	}
}
