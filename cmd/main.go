package main

import (
	"context"
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
	createUserHandler := user.PostUserHandler(userRepo)
	getUserHandler := user.GetUserByIdHandler(userRepo)

	mux.HandleFunc("POST /user/post/", createUserHandler)
	mux.HandleFunc("GET /user/get/", getUserHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
