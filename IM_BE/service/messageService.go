package service

import (
	"IM_BE/dto"
	"IM_BE/model"
	"IM_BE/repository"
	"IM_BE/utils"
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type MessageService struct {
	repo      repository.MessageRepo
	redisRepo repository.RedisRepo
}

func NewMessageService(repo repository.MessageRepo, redisRepo repository.RedisRepo) *MessageService {
	return &MessageService{repo: repo, redisRepo: redisRepo}
}

func (m *MessageService) SaveMessage(ctx context.Context, messageDto *dto.MessageRespDto) error {
	switch messageDto.ChatType {
	case dto.PrivateChat:
		message := &model.PrivateMessage{}
		if err := copier.Copy(message, messageDto); err != nil {
			utils.GetLogger().Error("复制消息体失败", zap.Error(err))
			return err
		}
		message.ReceiverId = messageDto.ToId

		seq, err := m.redisRepo.GetNextPrivateMessageSeq(ctx, message.SenderId, message.ReceiverId)
		if err != nil {
			utils.GetLogger().Error("redis 获取 private message seq 失败", zap.Error(err))
			return err
		}

		message.Seq = seq

		return m.repo.SavePrivateMessage(message)

	case dto.GroupChat:
		message := &model.GroupMessage{}
		if err := copier.Copy(message, messageDto); err != nil {
			utils.GetLogger().Error("复制群聊消息体失败", zap.Error(err))
			return err
		}
		message.GroupId = messageDto.ToId

		seq, err := m.redisRepo.GetNextGroupMessageSeq(ctx, message.GroupId)
		if err != nil {
			utils.GetLogger().Error("redis 获取 group message seq 失败", zap.Error(err))
			return err
		}

		message.Seq = seq

		return m.repo.SaveGroupMessage(message)

	default:
		utils.GetLogger().Error("chatType 错误", zap.Any("chatType:", messageDto.ChatType))
		return errors.New("chatType 错误")
	}
}

func (m *MessageService) GetHistoryMessages(ctx context.Context, senderId uint64, toId uint64, chatType dto.ChatType) ([]*dto.MessageRespDto, error) {
	switch chatType {
	case dto.PrivateChat:
		messages, err := m.repo.GetHistoryPrivateMessages(senderId, toId)
		if err != nil {
			return nil, err
		}

		var messageDtos []*dto.MessageRespDto
		for _, msg := range messages {
			messageDto := dto.MessageRespDto{
				ToId:     msg.ReceiverId,
				ChatType: "private",
			}
			err := copier.Copy(&messageDto, msg)
			if err != nil {
				utils.GetLogger().Error("复制消息体失败", zap.Error(err))
				continue
			}
			messageDtos = append(messageDtos, &messageDto)
		}

		return messageDtos, nil

	case dto.GroupChat:
		messages, err := m.repo.GetHistoryGroupMessages(toId)
		if err != nil {
			return nil, err
		}

		var messageDtos []*dto.MessageRespDto
		for _, msg := range messages {
			messageDto := dto.MessageRespDto{
				ToId:     msg.GroupId,
				ChatType: "group",
			}
			err := copier.Copy(&messageDto, msg)
			if err != nil {
				utils.GetLogger().Error("复制群聊消息体失败", zap.Error(err))
				continue
			}
			messageDtos = append(messageDtos, &messageDto)
		}

		return messageDtos, nil

	default:
		utils.GetLogger().Error("chatType 错误", zap.Any("chatType:", chatType))
		return nil, errors.New("chatType 错误")
	}
}
