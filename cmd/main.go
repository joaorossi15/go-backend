package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaorossi15/gobh/internal/user"
)

func main() {
	mux := http.NewServeMux()

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	userRepo := user.CreateUserRepo(pool)

	mux.HandleFunc("POST /user/post/", func(w http.ResponseWriter, r *http.Request) {
		id, err := userRepo.Create(context.Background(), "name", []byte("password"))
		if err != nil {
			log.Fatalf("%s", err)
		}

		fmt.Fprintf(w, "User %d created!", id)
	})

	mux.HandleFunc("GET /user/get/", func(w http.ResponseWriter, r *http.Request) {
		usr, err := userRepo.Get(context.Background(), "name")
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Fprintf(w, "User info: %d, %s", usr.ID, usr.Username)
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
