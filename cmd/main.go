package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaorossi15/gobh/internal/chat"
	"github.com/joaorossi15/gobh/internal/middleware"
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
	userLoginHandler := user.UserLoginHandler(userRepo)

	hub := chat.NewHub()
	go hub.Run(context.Background())
	chatHandler := chat.ChatHandler(hub, userRepo)

	mux.HandleFunc("POST /user/post/", createUserHandler)
	mux.HandleFunc("GET /user/get/", middleware.AuthMiddleware(getUserHandler))
	mux.HandleFunc("POST /user/login/", userLoginHandler)

	mux.HandleFunc("/ws/{roomID}", middleware.AuthMiddleware(middleware.NameMiddleware(chatHandler)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
