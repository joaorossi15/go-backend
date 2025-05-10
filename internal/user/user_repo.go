package repo

import (
	"context"

	"github.com/joaorossi15/gobh/internal/sqlc"
)

type UserRepo interface {
	GetUser(ctx context.Context, user string)
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams)
}


type UserR
