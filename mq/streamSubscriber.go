package mq

import (
	"IM_BE/dto"
	"IM_BE/utils"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type StreamSubscriber struct {
	rdb       *redis.Client
	stream    string
	group     string
	consumer  string
	onMessage func(*dto.MessageRespDto) error

	// 对应的 client 离线之后，消费者组需要删除，这里就体现出来 ctx 的中断的作用
	ctx    context.Context
	cancel context.CancelFunc
}

func NewStreamSubscriber(rdb *redis.Client, stream string, group string, consumer string, handler func(*dto.MessageRespDto) error) *StreamSubscriber {
	ctx, cancel := context.WithCancel(context.Background())

	return &StreamSubscriber{
		rdb:       rdb,
		stream:    stream,
		group:     group,
		consumer:  consumer,
		onMessage: handler,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (s *StreamSubscriber) Cancel() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *StreamSubscriber) InitGroup() error {
	_, err := s.rdb.XGroupCreateMkStream(context.Background(), s.stream, s.group, "$").Result() // "$" 从 Stream 尾部开始，只消费新写入的消息
	if err != nil {
		utils.GetLogger().Error("创建消费者组失败", zap.Error(err))
		return err
	}
	return nil
}

// 开始消费
func (s *StreamSubscriber) Start() {
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				s.consume()
			}
		}
	}()
}

func (s *StreamSubscriber) consume() {
	utils.GetLogger().Debug("开始消费")

	res, err := s.rdb.XReadGroup(context.Background(), &redis.XReadGroupArgs{
		Group:    s.group,
		Consumer: s.consumer,
		Streams:  []string{s.stream, ">"},
		Block:    5 * time.Second,
		Count:    1,
	}).Result()

	if err != nil {
		if err.Error() == "redis: nil" {
			utils.GetLogger().Debug("暂无新消息")
			return
		}
		utils.GetLogger().Error("读取消息失败", zap.Error(err))
		return
	}

	for _, stream := range res {
		for _, msg := range stream.Messages {
			var rawMessage *dto.RawMessage
			jsonData, _ := json.Marshal(msg.Values)
			if err := json.Unmarshal(jsonData, &rawMessage); err != nil {
				utils.GetLogger().Error("消息解析失败", zap.Error(err), zap.Any("msg", msg))
				break
			}

			message, err := rawMessage.ToMessageRespDto()
			if err != nil {
				utils.GetLogger().Error("消息类型转换失败", zap.Error(err))
				break
			}

			// 执行回调（将消息放进对应的信箱）
			if s.onMessage != nil && message != nil {
				if err := s.onMessage(message); err != nil {
					utils.GetLogger().Error("消息回调执行失败")
					break
				}
			}

			// ack
			s.rdb.XAck(context.Background(), s.stream, s.group, msg.ID)
			utils.GetLogger().Debug("消息 ack")
		}
	}
}
