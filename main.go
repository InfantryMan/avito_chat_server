package main

import (
	"github.com/InfantryMan/avito_chat_server/app"
	"github.com/InfantryMan/avito_chat_server/databases"
	_ "github.com/InfantryMan/avito_chat_server/models"
	"github.com/InfantryMan/avito_chat_server/routes"
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
	router.POST("/chats/get", routes.GetChats)
	router.POST("/messages/get", routes.GetMessages)

	host, port := app.GetHostAndPort()
	addr := host + ":" + port
	log.Println("Application is running on " + addr)

	log.Fatal(http.ListenAndServe(addr, router))
}
