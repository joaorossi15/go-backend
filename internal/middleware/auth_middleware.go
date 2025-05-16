package middleware

import (
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk := r.Header.Get("Authorization")

		if tk == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(tk, " ", 2)
		if err := VerifyToken(parts[1]); err != nil {
			http.Error(w, "invalid jwt token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
