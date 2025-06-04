package iface

import (
	"IM_BE/dto"
	"context"
)

// MessageHandler 解决循环导入
type MessageHandler interface {
	SaveMessage(ctx context.Context, message *dto.MessageRespDto) error
	GetHistoryMessages(ctx context.Context, senderId uint64, toId uint64, chatType dto.ChatType) ([]*dto.MessageRespDto, error)
}
