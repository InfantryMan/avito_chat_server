package main

import (
	"awesomeProject/databases"
	_ "awesomeProject/models"
	"awesomeProject/routes"
	"github.com/dimfeld/httptreemux"
	"log"
	"net/http"
)

func main() {
	databases.ConnectDB()
	defer databases.CloseDB()

	router := httptreemux.New()
	router.POST("/users/add", routes.AddUser)
	router.POST("/chats/add", routes.AddChat)
	router.POST("/messages/add", routes.AddMessage)

	log.Fatal(http.ListenAndServe(":9000", router))
}
