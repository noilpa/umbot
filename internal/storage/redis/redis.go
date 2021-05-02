package redis

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis/v8"

	"github.com/noilpa/umbot/internal/app"
)

type client struct {
	db *redis.Client
}

func New(c Config) *client {
	return &client{
		db: redis.NewClient(&redis.Options{
			Addr:     c.Address,
			Username: c.Login,
			Password: c.Password,
		}),
	}
}

func (c *client) Get(ctx context.Context, userID int) (app.UserData, error) {
	var out app.UserData

	res, err := c.db.Get(ctx, strconv.Itoa(userID)).Bytes()
	if err != nil {
		return out, err
	}

	if err := json.Unmarshal(res, &out); err != nil {
		return out, err
	}

	return out, nil
}

func (c *client) Set(ctx context.Context, userID int, data app.UserData) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.db.Set(ctx, strconv.Itoa(userID), value, 0).Err()
}
