package redis

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis/v8"

	"github.com/noilpa/umbot/internal/storage/entity"
)

type client struct {
	db *redis.Client
}

func New(c Config) *client {
	db := redis.NewClient(&redis.Options{
		Addr:     c.Address,
		Username: c.Login,
		Password: c.Password,
	})

	if err := db.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &client{
		db: db,
	}
}

func (c *client) Get(ctx context.Context, chatID int64) (entity.ChatData, error) {
	var out entity.ChatData

	res, err := c.db.Get(ctx, strconv.FormatInt(chatID, 10)).Bytes()
	if err != nil {
		return out, err
	}

	if err := json.Unmarshal(res, &out); err != nil {
		return out, err
	}

	return out, nil
}

func (c *client) Set(ctx context.Context, chatID int64, data entity.ChatData) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.db.Set(ctx, strconv.FormatInt(chatID,10), value, 0).Err()
}
