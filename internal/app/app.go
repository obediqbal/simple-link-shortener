package app

import (
	"fmt"
	"log"
	"net/http"
	"simplelinkshortener/internal/app/handler"

	"github.com/gorilla/mux"
)

func Init() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler.HomeHandler).Methods("GET")
	r.HandleFunc("/shorten", handler.ShortenHandler).Methods("POST")
	r.HandleFunc("/{shortURL}", handler.RedirectHandler).Methods("GET")

	port := ":8080"
	fmt.Printf("Server is listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
