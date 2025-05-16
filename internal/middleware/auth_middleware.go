package middleware

import (
	"context"
	"net/http"
	"strings"
)

type usrKey string

const userKey usrKey = "userid"

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk := r.Header.Get("Authorization")

		if tk == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(tk, " ", 2)
		if _, err := VerifyToken(parts[1]); err != nil {
			http.Error(w, "invalid jwt token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, parts[1])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
