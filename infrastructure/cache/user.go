package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user-service/application/cache"
	"user-service/domain/entities"
	"user-service/domain/errs"
	db "user-service/infrastructure"

	"github.com/redis/go-redis/v9"
)

type userCache struct {
	rdb *redis.Client
	ctx context.Context
}

func NewUserCache(redisURL string) cache.UserCache {
	return &userCache{
		rdb: db.NewRedisClient(redisURL),
		ctx: context.Background(),
	}
}

func handleRedisError(err error) error {
	switch (err) {
	case redis.Nil:
		return errs.ErrUserNotFound
	default:
		return errs.ErrCaching
	}
}

func (c *userCache) SetUser(userID uint, user entities.User) error {
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		return handleRedisError(err)
	}

	key := fmt.Sprintf("user:%d", userID)
	return c.rdb.Set(c.ctx, key, string(jsonBytes), time.Minute * 30).Err()
}

func (c *userCache) GetUser(userID uint) (entities.User, error) {
	val, err := c.rdb.Get(c.ctx, fmt.Sprintf("user:%d", userID)).Result()
	if err != nil {
		return entities.User{}, err
	}

	var user entities.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return entities.User{}, err
	}

	return user, nil
}
