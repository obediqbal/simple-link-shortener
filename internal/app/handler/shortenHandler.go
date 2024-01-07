package handler

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"simplelinkshortener/internal/pkg/database"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("shortening...")
	host := os.Getenv("APP_HOST")

	originalURL := r.FormValue("original_url")
	if originalURL == "" {
		http.Error(w, "Please provide a URL", http.StatusBadRequest)
		return
	}
	u, err := url.Parse(originalURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := publicsuffix.EffectiveTLDPlusOne(u.Host); err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	db := database.New()
	var shortURL string
	if shortURL = r.FormValue("short_url"); shortURL == "" { // Generates short_url if not provided
		for {
			shortURL = generateShortURL()
			var exists bool
			if err := isShortURLExist(shortURL, &exists); err != nil {
				log.Println("Error validating short URL ", err)
				http.Error(w, "Error validating short URL", http.StatusInternalServerError)
				return
			}
			if !exists {
				break
			}
		}
	} else {
		var exists bool
		if err := isShortURLExist(shortURL, &exists); err != nil {
			log.Println("Error validating short URL ", err)
			http.Error(w, "Error validating short URL", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Short URL already used", http.StatusBadRequest)
			return
		}
	}

	if _, err := url.ParseRequestURI(host + "/" + shortURL); err != nil { // Validates the short_url provided or generated
		log.Println("Invalid short URL ", err)
		http.Error(w, "Invalid short URL", http.StatusBadRequest)
		return
	}

	if _, err := db.Exec("INSERT INTO url(original_url, short_url) VALUES($1, $2)", originalURL, shortURL); err != nil {
		log.Println("Error inserting data ", err)
		http.Error(w, "Error inserting data", http.StatusInternalServerError)
		return
	}
	// shortenedURLs[shortURL] = originalURL

	response := fmt.Sprintf("Shortened URL: %s/%s", host, shortURL)
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

func isShortURLExist(shortURL string, dest *bool) error {
	db := database.New()
	row := db.QueryRow("SELECT short_url FROM url WHERE short_url=$1", shortURL)
	if err := row.Scan(); err != nil {
		if err == sql.ErrNoRows {
			*dest = false
		} else {
			return err
		}
	}
	*dest = true
	return nil
}
