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
	"strconv"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("shortening...")
	host := os.Getenv("APP_TARGET")

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
		if max_len, _ := strconv.Atoi(os.Getenv("MAX_SHORT_URL_LEN")); len(shortURL) > max_len {
			http.Error(w, fmt.Sprintf("Short URL is too long (max %d)", max_len), http.StatusBadRequest)
			return
		}
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

	finalShortURL, err := url.ParseRequestURI(host + "/" + shortURL)
	if err != nil { // Validates the short_url provided or generated
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

	response := finalShortURL.String()
	fmt.Fprint(w, response)
	fmt.Println(response)
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
	var data string
	if err := row.Scan(&data); err != nil {
		if err == sql.ErrNoRows {
			*dest = false
			return nil
		} else {
			return err
		}
	}
	*dest = true
	return nil
}
