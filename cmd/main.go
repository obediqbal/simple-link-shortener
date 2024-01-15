package main

import (
	"fmt"
	"log"
	"simplelinkshortener/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed loading .env")
	} else {
		log.Println("Loaded .env file")
	}
	fmt.Printf("Starting app!!7!\n")
	app.Init()
}
