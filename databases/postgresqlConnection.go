package databases

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

var (
	host     = "localhost"
	port     = "5432"
	user     = "chat_admin"
	password = "chat_password"
	dbname   = "db_chat"
)

func ConnectDB() {
	e := godotenv.Load()
	if e == nil {
		user = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASS")
		dbname = os.Getenv("DB_NAME")
		host = os.Getenv("DB_HOST")
		port = os.Getenv("DB_PORT")
	} else {
		log.Println("ConnectDB: .env file not found.\n" +
			"Default settings are used.")
	}

	dbUri := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	fmt.Println(dbUri)

	conn, err := sql.Open("postgres", dbUri)
	if err != nil {
		log.Fatalln("ConnectDB: unable to connect to DB\n", err)
	}
	db = conn

	err = db.Ping()
	if err != nil {
		log.Fatalln("ConnectDB: unable to ping DB\n", err)
	}

	log.Println("Successfully connected to database " + dbname)
}

func GetPostgresSession() *sql.DB {
	return db
}

func CloseDB() {
	db.Close()
}
