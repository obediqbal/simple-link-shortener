package main

import (
	"fmt"
	"log"
	"os"
	"simplelinkshortener/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) > 2 && os.Args[2] == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file", err)
		} else {
			log.Println("Loaded .env file")
		}
	}
	fmt.Printf("Starting app!!7!\n")
	app.Init()
}
