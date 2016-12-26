package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mattmac4241/chat-auth/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbname := os.Getenv("DBNAME")
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := service.InitDatabase(dbinfo)
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	service.DB = db
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	server := service.NewServer()
	server.Run(":" + port)
}
