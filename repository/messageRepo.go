package repository

import (
	"IM_BE/model"
	"IM_BE/utils"
	"database/sql"
	"go.uber.org/zap"
)

type MessageRepo interface {
	SavePrivateMessage(*model.PrivateMessage) error // TODO:参数不要轻易使用指针
	SaveGroupMessage(*model.GroupMessage) error
	GetHistoryPrivateMessages(senderId uint64, toId uint64) ([]*model.PrivateMessage, error)
	GetHistoryGroupMessages(groupId uint64) ([]*model.GroupMessage, error)
}

func NewMessageRepo(db *sql.DB) MessageRepo {
	return &messageRepoImpl{db: db}
}

type messageRepoImpl struct {
	db *sql.DB
}

func (m *messageRepoImpl) GetHistoryPrivateMessages(senderId uint64, toId uint64) ([]*model.PrivateMessage, error) {
	query := "select id, sender_id, content_type, content, created_at from private_messages where status = '0' and sender_id = ? and receiver_id = ? order by seq desc limit 10"

	rows, err := m.db.Query(query, senderId, toId)
	if err != nil {
		utils.GetLogger().Error("查询失败", zap.Error(err))
		return nil, err
	}

	var messages []*model.PrivateMessage
	for rows.Next() {
		message := &model.PrivateMessage{}

		if err = rows.Scan(&message.Id, &message.SenderId, &message.ContentType, &message.Content, &message.CreatedAt); err != nil {
			utils.GetLogger().Error("读取结果失败", zap.Error(err))
			return nil, err
		}

		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		utils.GetLogger().Error("遍历行失败", zap.Error(err))
		return nil, err
	}

	return messages, nil
}

func (m *messageRepoImpl) GetHistoryGroupMessages(groupId uint64) ([]*model.GroupMessage, error) {
	query := "select id, sender_id, content_type, content, created_at from group_messages where status = '0' and group_id = ? order by seq desc limit 10"

	rows, err := m.db.Query(query, groupId)
	if err != nil {
		utils.GetLogger().Error("查询失败", zap.Error(err))
		return nil, err
	}

	var messages []*model.GroupMessage
	for rows.Next() {
		message := &model.GroupMessage{}

		if err = rows.Scan(&message.Id, &message.SenderId, &message.ContentType, &message.Content, &message.CreatedAt); err != nil {
			utils.GetLogger().Error("读取结果失败", zap.Error(err))
			return nil, err
		}

		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		utils.GetLogger().Error("遍历行失败", zap.Error(err))
		return nil, err
	}

	return messages, nil
}

func (m *messageRepoImpl) SavePrivateMessage(message *model.PrivateMessage) error {
	query := "insert into private_messages(sender_id, receiver_id, content_type, content, seq) values(?, ?, ?, ?, ?)"

	if _, err := m.db.Exec(query, message.SenderId, message.ReceiverId, message.ContentType, message.Content, message.Seq); err != nil {
		utils.GetLogger().Error("插入失败", zap.Error(err))
		return err
	}

	return nil
}

func (m *messageRepoImpl) SaveGroupMessage(message *model.GroupMessage) error {
	query := "insert into group_messages(sender_id, group_id, content_type, content, seq) values(?, ?, ?, ?, ?)"

	if _, err := m.db.Exec(query, message.SenderId, message.GroupId, message.ContentType, message.Content, message.Seq); err != nil {
		utils.GetLogger().Error("插入失败", zap.Error(err))
		return err
	}

	return nil
}
