package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	shortenedURLs map[string]string
	letters       = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
)

func init() {
	shortenedURLs = make(map[string]string)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/shorten", shortenHandler).Methods("POST")
	r.HandleFunc("/{shortURL}", redirectHandler).Methods("GET")

	port := ":8080"
	fmt.Printf("Server is listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world!!!")
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("original_url")
	if originalURL == "" {
		http.Error(w, "Please prove a URL", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()
	shortenedURLs[shortURL] = originalURL

	response := fmt.Sprintf("Shortened URL: http://localhost:8080/%s", shortURL)
	fmt.Fprintf(w, response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	originalURL, exists := shortenedURLs[shortURL]
	if !exists {
		http.Error(w, "Shortened URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func generateShortURL() string {
	var builder strings.Builder
	for i := 0; i < 6; i++ {
		builder.WriteByte(letters[rand.Intn(len(letters))])
	}
	return builder.String()
}
