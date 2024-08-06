package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gindemo/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

//go:generate mockgen -source=./user.go -package=cachemocks -destination=./mocks/user.mock.go UserCache
type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
	Del(ctx context.Context, id int64) error
}

type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (c *RedisUserCache) Del(ctx context.Context, id int64) error {
	return c.cmd.Del(ctx, c.key(id)).Err()
}

func (c *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.key(uid)
	// 暂时用JSON
	data, err := c.cmd.Get(ctx, key).Result()
	//data, err := c.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	//if err != nil {
	//	return domain.User{}, err
	//}
	return u, err
}

func (c *RedisUserCache) Set(ctx context.Context, du domain.User) error {
	key := c.key(du.Id)
	// 暂时用JSON
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisUserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}

type UserCacheV1 struct {
	client *redis.Client
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
