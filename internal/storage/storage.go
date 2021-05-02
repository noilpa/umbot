package storage

import (
	"context"

	"github.com/noilpa/umbot/internal/app"
)

type Storage interface {
	Get(ctx context.Context, userID int) (app.UserData, error)
	Set(ctx context.Context, userID int, data app.UserData) error
}
