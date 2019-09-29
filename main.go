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

	log.Fatal(http.ListenAndServe(":9000", router))
}
