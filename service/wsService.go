package service

import (
	"IM_BE/db/mysql"
	"IM_BE/db/redis"
	"IM_BE/dto"
	"IM_BE/mq"
	"IM_BE/repository"
	"IM_BE/utils"
	"IM_BE/ws"
	"context"
	"errors"
	"fmt"
)

type WsService struct {
	repo repository.WsRepo
}

func NewWsService(repo repository.WsRepo) *WsService {
	return &WsService{repo: repo}
}

func (w *WsService) AddClient(ctx context.Context, client *ws.Client) error {
	manager := ws.GetWsManager()

	// TODO：并发安全，底层原理
	manager.Register <- client

	userId, ok := ctx.Value("user_id").(uint64)
	if !ok {
		utils.GetLogger().Error("类型断言失败")
		return errors.New("类型断言失败")
	}

	// 初始化对应的消费者组
	messageService := NewMessageService(
		repository.NewMessageRepo(mysql.Get()),
		repository.NewRedisRepo(redis.Get()),
	)
	// 辅助函数：创建并启动 StreamSubscriber
	startSubscriber := func(stream, group, consumer string) error {
		subscriber := mq.NewStreamSubscriber(redis.Get(), stream, group, consumer, func(msg *dto.MessageRespDto) error {
			// 持久化消息
			if err := messageService.SaveMessage(context.Background(), msg); err != nil {
				return err
			}
			utils.GetLogger().Debug("持久化消息")

			// 将消息放进信箱
			client.GetMessage(msg)
			utils.GetLogger().Debug("将消息放进信箱")

			return nil
		})

		if err := subscriber.InitGroup(); err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
			return err
		}

		client.AddSubscriber(subscriber)
		subscriber.Start()
		return nil
	}

	// 添加好友订阅（私聊）
	friendIds, err := w.repo.GetUserFriendShips(ctx)
	if err != nil {
		return err
	}
	for _, fid := range friendIds {
		user1, user2 := userId, fid
		if user2 < user1 {
			user1, user2 = user2, user1
		}
		stream := fmt.Sprintf("chat:user:%d:%d", user1, user2)
		if err := startSubscriber(stream, fmt.Sprintf("group:%d", userId), fmt.Sprintf("consumer:%d", userId)); err != nil {
			return err
		}
	}
	utils.GetLogger().Debug("好友加载完毕")

	// 添加群聊订阅
	groupIds, err := w.repo.GetUserGroupShips(ctx)
	if err != nil {
		return err
	}
	for _, gid := range groupIds {
		stream := fmt.Sprintf("chat:group:%d", gid)
		if err := startSubscriber(stream, fmt.Sprintf("group:%d", userId), fmt.Sprintf("consumer:%d", userId)); err != nil {
			return err
		}
	}
	utils.GetLogger().Debug("群聊加载完毕")

	// TODO：协程的内存泄露
	// ws IO
	go client.WritePump()
	go client.ReadPump()

	return nil
}
