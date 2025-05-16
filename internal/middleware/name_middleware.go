package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

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
		auth := r.Header.Get("Authorization")

		if auth == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(auth, "Bearer ")
		auth = splitToken[1]

		token, err := jwt.ParseWithClaims(auth, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			http.Error(w, "invalid jwt token", http.StatusUnauthorized)
			return
		}

		claims, _ := token.Claims.(*MyCustomClaims)

		ctx := context.WithValue(r.Context(), UserIDKey, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
