package storage

import (
	"context"

	"github.com/noilpa/umbot/internal/storage/entity"
)

type Storage interface {
	Get(ctx context.Context, chatID int64) (entity.ChatData, error)
	Set(ctx context.Context, chatID int64, data entity.ChatData) error
}
