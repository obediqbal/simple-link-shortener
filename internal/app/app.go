package app

import (
	"fmt"
	"log"
	"net/http"
	"simplelinkshortener/internal/app/handler"
	"simplelinkshortener/internal/app/middlewares"

	"github.com/gorilla/mux"
)

func Init() {
	r := mux.NewRouter()

	// allowedOrigins := strings.Split(os.Getenv("ORIGIN_ALLOWED"), ",")
	// log.Printf("Allowed: %v\n", allowedOrigins)

	// corsHandler := handlers.CORS(
	// 	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	// 	handlers.AllowedMethods([]string{"GET", "POST"}),
	// 	handlers.AllowedOrigins(allowedOrigins),
	// )

	// r.Use(corsHandler)
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
	fmt.Printf("Server is listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
