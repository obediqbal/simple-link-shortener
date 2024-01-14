package middlewares

import (
	"context"
	"net/http"
	"simplelinkshortener/internal/pkg/auth"
)

func UseAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenString := r.Header.Get("Authorization"); tokenString != "" {
			claims, err := auth.ValidateJWT(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), "username", claims.Username)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}
	})
}

func UseLazyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := ""
		if tokenString := r.Header.Get("Authorization"); tokenString != "" {
			claims, err := auth.ValidateJWT(tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			username = claims.Username
		}
		ctx := context.WithValue(r.Context(), "username", username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
