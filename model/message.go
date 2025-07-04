package model

import (
	"time"
)

type MessageType string

const (
	TextMessage   MessageType = "text"
	ImageMessage  MessageType = "image"
	FileMessage   MessageType = "file"
	SystemNotice  MessageType = "system"
	FriendRequest MessageType = "friendRequest"
	AudioMessage  MessageType = "audio"
)

// copier只能实现非匿名字段的复制
type messageBasic struct {
	Id          uint64
	SenderId    uint64
	ContentType MessageType
	Content     string
	CreatedAt   time.Time
	Seq         uint64
}

type PrivateMessage struct {
	Id          uint64
	SenderId    uint64
	ReceiverId  uint64
	ContentType MessageType
	Content     string
	CreatedAt   time.Time
	Seq         uint64 `json:"seq"`
}

type GroupMessage struct {
	Id          uint64
	SenderId    uint64
	GroupId     uint64
	ContentType MessageType
	Content     string
	CreatedAt   time.Time
	Seq         uint64 `json:"seq"`
}
