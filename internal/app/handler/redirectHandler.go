package handler

import (
	"fmt"
	"log"
	"net/http"
	"simplelinkshortener/internal/pkg/database"

	"github.com/gorilla/mux"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]
	if shortURL == "" || shortURL == "favicon.ico" {
		return
	}
	fmt.Printf("redirecting...%s", shortURL)

	db := database.New()
	row := db.QueryRow("SELECT original_url FROM url WHERE short_url=$1", shortURL)

	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		log.Println("Shortened URL not found", err)
		http.Error(w, "Shortened URL not found", http.StatusNotFound)
		return
	}

	// originalURL, exists := shortenedURLs[shortURL]
	// if !exists {
	// 	return
	// }

	fmt.Printf("Redirecting to %v", originalURL)
	http.Redirect(w, r, originalURL, http.StatusFound)
}
