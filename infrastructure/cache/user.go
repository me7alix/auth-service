package cache

import (
	"fmt"
	"time"
	"errors"
	"context"
	"encoding/json"
	"user-service/application/cache"
	"user-service/domain/entities"
	"user-service/domain/errs"
	db "user-service/infrastructure"

	"github.com/phuslu/log"
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
	if err == nil { return nil }
	if errors.Is(err, redis.Nil) {
		return errs.ErrUserNotFound
	} else {
		return errs.ErrCaching
	}
}

func (c *userCache) SetUser(userID uint, user entities.User) error {
	log.Debug().Msgf("cache: SetUser(%v) %v", userID, user)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		log.Debug().Msg(err.Error())
		return errs.ErrCaching
	}
	key := fmt.Sprintf("user:%d", userID)
	err = c.rdb.Set(c.ctx, key, string(jsonBytes), time.Minute * 30).Err()
	return handleRedisError(err)
}

func (c *userCache) GetUser(userID uint) (entities.User, error) {
	val, err := c.rdb.Get(c.ctx, fmt.Sprintf("user:%d", userID)).Result()
	if err != nil {
		return entities.User{}, handleRedisError(err)
	}
	var user entities.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return entities.User{}, handleRedisError(err)
	}
	log.Debug().Msgf("cache: GetUser(%v) %v", userID, user)
	return user, nil
}
