package app

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	port = "9000"
	host = "localhost"
)

func GetHostAndPort() (string, string) {
	e := godotenv.Load()
	if e == nil {
		host = os.Getenv("APP_HOST")
		port = os.Getenv("APP_PORT")
	} else {
		log.Println("GetHostAndPort: .env file not found.\n" +
			"Default settings are used.")
	}
	return host, port
}
