package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// TODO: redis 是单线程的，所有操作均有原子性，都是安全的

type RedisRepo interface {
	SetUserToken(ctx context.Context, userId uint64, token string, expiration int) error
	GetUserToken(ctx context.Context, userId uint64) (string, error)
	DelUserToken(ctx context.Context, userId uint64) error

	GetNextPrivateMessageSeq(ctx context.Context, user1 uint64, user2 uint64) (uint64, error)
	GetNextGroupMessageSeq(ctx context.Context, groupId uint64) (uint64, error)
}

func NewRedisRepo(rdb *redis.Client) RedisRepo {
	return &redisRepoImpl{rdb: rdb}
}

type redisRepoImpl struct {
	rdb *redis.Client
}

// TODO:rdb的方法的方法？？？？

func (c *redisRepoImpl) SetUserToken(ctx context.Context, userId uint64, token string, expiration int) error {
	tokenKey := fmt.Sprintf("user_%d_token", userId)
	// 设置时间，自动清除 token，不占内存
	return c.rdb.Set(ctx, tokenKey, token, time.Duration(expiration)*time.Hour).Err()
}

func (c *redisRepoImpl) GetUserToken(ctx context.Context, userId uint64) (string, error) {
	tokenKey := fmt.Sprintf("user_%d_token", userId)
	return c.rdb.Get(ctx, tokenKey).Result()
}

func (c *redisRepoImpl) DelUserToken(ctx context.Context, userId uint64) error {
	tokenKey := fmt.Sprintf("user_%d_token", userId)
	return c.rdb.Del(ctx, tokenKey).Err()
}

func (c *redisRepoImpl) GetNextPrivateMessageSeq(ctx context.Context, user1 uint64, user2 uint64) (uint64, error) {
	if user1 < user2 {
		user1, user2 = user2, user1 // TODO:go中的赋值号
	}

	seqKey := fmt.Sprintf("private_%d_%d_seq", user1, user2)
	return c.rdb.Incr(ctx, seqKey).Uint64()
}

func (c *redisRepoImpl) GetNextGroupMessageSeq(ctx context.Context, groupId uint64) (uint64, error) {
	seqKey := fmt.Sprintf("group_%d_seq", groupId)
	return c.rdb.Incr(ctx, seqKey).Uint64()
}
