package app

import (
	"log"
	"net/http"
	"simplelinkshortener/internal/app/handler"
	"simplelinkshortener/internal/app/middlewares"

	"github.com/gorilla/mux"
)

func Init() {
	r := mux.NewRouter()

	r.Use(middlewares.UseCors)

	publicRouter := r.NewRoute().Subrouter()
	publicRouter.HandleFunc("/", handler.HomeHandler).Methods("GET")
	publicRouter.HandleFunc("/auth/register", handler.RegisterHandler).Methods("POST")
	publicRouter.HandleFunc("/auth/login", handler.LoginHandler).Methods("POST")
	publicRouter.HandleFunc("/{shortURL}", handler.RedirectHandler).Methods("GET")

	secureRouter := r.NewRoute().Subrouter()
	secureRouter.Use(middlewares.UseAuth)

	lazyRouter := r.NewRoute().Subrouter()
	lazyRouter.Use(middlewares.UseLazyAuth)
	lazyRouter.HandleFunc("/shorten", handler.ShortenHandler).Methods("POST")

	port := ":8080"
	log.Printf("Server is listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
