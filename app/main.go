package main

import (
	"log"

	"kintai-bot/app/infrastructure"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Could not load the .env file.")
	}

	infrastructure.Init()
}
