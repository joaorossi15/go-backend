package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaorossi15/gobh/internal/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Get(ctx context.Context, user string)
	Create(ctx context.Context)
}

type UserR struct {
	q *sqlc.Queries
}

func CreateUserRepo(pool *pgxpool.Pool) *UserR {
	return &UserR{q: sqlc.New(pool)}
}

func (usr *UserR) Get(ctx context.Context, name string) (sqlc.User, error) {
	user, err := usr.q.GetUser(ctx, name)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}

func (usr *UserR) Create(ctx context.Context, name string, password []byte) (int64, error) {
	pwd, err := bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		return 0, err
	}

	usrRow, err := usr.q.CreateUser(ctx, sqlc.CreateUserParams{
		Username: name,
		Password: pwd,
	})
	if err != nil {
		return 0, err
	}

	return usrRow.ID, nil
}
