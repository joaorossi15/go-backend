package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const UserIDKey ctxKey = "uid"

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NameMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(userKey)
		tk, ok := v.(string)

		if !ok {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		token, err := VerifyToken(tk)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, _ := token.Claims.(*MyCustomClaims)

		ctx := context.WithValue(r.Context(), UserIDKey, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
