package dto

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	TextMessage  MessageType = "text"
	ImageMessage MessageType = "image"
	FileMessage  MessageType = "file"
	SystemNotice MessageType = "system"
	AudioMessage MessageType = "audio"
)

type ChatType string

const (
	PrivateChat ChatType = "private"
	GroupChar   ChatType = "group"
)

// TODO:怎么避免使用id区分不同的对象
type MessageDto struct {
	SenderId    uint64          `json:"sender_id,omitempty"`
	ToId        uint64          `json:"to_id,omitempty"`
	ChatType    ChatType        `json:"chat_type"`
	ContentType MessageType     `json:"content_type"` // TODO:不必担心前端是stringer导致json解析失败
	Content     json.RawMessage `json:"content"`
	CreatedAt   time.Time       `json:"created_at"`
}
