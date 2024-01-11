package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"simplelinkshortener/internal/app/handler"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Init() {
	r := mux.NewRouter()

	allowedOrigins := strings.Split(os.Getenv("ORIGIN_ALLOWED"), ",")
	log.Printf("Allowed: %v\n", allowedOrigins)

	corsHandler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedOrigins(allowedOrigins),
	)

	r.HandleFunc("/", handler.HomeHandler).Methods("GET")
	r.HandleFunc("/shorten", handler.ShortenHandler).Methods("POST")
	r.HandleFunc("/{shortURL}", handler.RedirectHandler).Methods("GET")

	r.Use(corsHandler)

	port := ":8080"
	fmt.Printf("Server is listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
