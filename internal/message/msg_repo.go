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

func (m *MessageRepo) GetConversationMessages(ctx context.Context, senderID int64, recID int64) ([]sqlc.Message, error) {
	messages, err := m.q.SelectConvMessages(ctx, sqlc.SelectConvMessagesParams{
		SenderID: senderID,
		RecID:    recID,
	})

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *MessageRepo) GetUserSentMessages(ctx context.Context, userID int64) ([]sqlc.Message, error) {
	messages, err := m.q.SelectSentMessages(ctx, userID)

	if err != nil {
		return nil, err
	}

	return messages, nil

}

func (m *MessageRepo) GetUserReceivedMessages(ctx context.Context, userID int64) ([]sqlc.Message, error) {
	messages, err := m.q.SelectRecMessages(ctx, userID)

	if err != nil {
		return nil, err
	}

	return messages, nil

}
