package message

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaorossi15/gobh/internal/sqlc"
)

type MessageRepo struct {
	q *sqlc.Queries
}

func CreateMessageRepo(pool *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{q: sqlc.New(pool)}
}

func (m *MessageRepo) CreateMessage(ctx context.Context, senderID int64, recID int64, body string) (int64, error) {
	message, err := m.q.CreateMessage(ctx, sqlc.CreateMessageParams{
		SenderID: senderID,
		RecID:    recID,
		Body:     body,
	})

	if err != nil {
		return 0, err
	}

	return message.ID, nil
}

