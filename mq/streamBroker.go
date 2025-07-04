package mq

import (
	"IM_BE/dto"
	"context"
	"github.com/redis/go-redis/v9"
)

type StreamBroker struct {
	rdb *redis.Client
}

func NewStreamBroker(rdb *redis.Client) *StreamBroker {
	return &StreamBroker{rdb: rdb}
}

func (s *StreamBroker) Publish(ctx context.Context, stream string, msg *dto.MessageRespDto) (string, error) { // 返回消息的 id
	return s.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: msg.ToMap(),
	}).Result()
}
