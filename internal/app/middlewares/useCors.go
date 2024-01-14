package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
)

func UseCors(next http.Handler) http.Handler {
	allowedOrigins := strings.Split(os.Getenv("ORIGIN_ALLOWED"), ",")
	handler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedOrigins(allowedOrigins),
	)

	return handler(next)
}
