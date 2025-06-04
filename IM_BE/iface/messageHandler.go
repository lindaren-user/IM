package iface

import "IM_BE/dto"

// MessageHandler 解决循环导入
type MessageHandler interface {
	SaveMessage(message *dto.MessageRespDto) error
	GetHistoryMessages(senderId uint64, toId uint64, chatType dto.ChatType) ([]*dto.MessageRespDto, error)
}
