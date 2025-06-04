package dto

import (
	"IM_BE/model"
	"encoding/json"
	"time"
)

type ChatType string

const (
	PrivateChat ChatType = "private"
	GroupChat   ChatType = "group"
)

// TODO: copier无法复制匿名字段
type messageBasic struct {
	SenderId    uint64            `json:"sender_id,omitempty"`
	ToId        uint64            `json:"to_id,omitempty"`
	ChatType    ChatType          `json:"chat_type"`
	ContentType model.MessageType `json:"content_type"` // TODO:不必担心前端是string而导致json解析失败
	CreatedAt   time.Time         `json:"created_at"`
}

type MessageReqDto struct {
	SenderId    uint64            `json:"sender_id,omitempty"`
	ToId        uint64            `json:"to_id,omitempty"`
	ChatType    ChatType          `json:"chat_type"`
	ContentType model.MessageType `json:"content_type"` // TODO:不必担心前端是string而导致json解析失败
	CreatedAt   time.Time         `json:"created_at"`
	Content     json.RawMessage   `json:"content"`
}

type MessageRespDto struct {
	Id          uint64            `json:"id"`
	SenderId    uint64            `json:"sender_id,omitempty"`
	ToId        uint64            `json:"to_id,omitempty"`
	ChatType    ChatType          `json:"chat_type"`
	ContentType model.MessageType `json:"content_type"` // TODO:不必担心前端是string而导致json解析失败
	CreatedAt   time.Time         `json:"created_at"`
	Content     string            `json:"content"`
}
