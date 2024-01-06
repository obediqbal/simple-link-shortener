package handler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"simplelinkshortener/internal/pkg/database"
	"strings"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("shortening...")

	originalURL := r.FormValue("original_url")
	if originalURL == "" {
		http.Error(w, "Please provide a URL", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()

	db := database.New()
	_, err := db.Exec("INSERT INTO url(original_url, short_url) VALUES($1, $2)", originalURL, shortURL)
	if err != nil {
		log.Println("Error inserting data", err)
		http.Error(w, "Error inserting data", http.StatusInternalServerError)
		return
	}
	// shortenedURLs[shortURL] = originalURL

	response := fmt.Sprintf("Shortened URL: http://localhost:8080/%s", shortURL)
	fmt.Fprintf(w, response)
}

func generateShortURL() string {
	letters := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	var builder strings.Builder
	for i := 0; i < 6; i++ {
		builder.WriteByte(letters[rand.Intn(len(letters))])
	}
	return builder.String()
}
