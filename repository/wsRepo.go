package repository

import (
	"IM_BE/utils"
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

// 名字以业务划分
type WsRepo interface {
	GetUserFriendShips(ctx context.Context) ([]uint64, error)
	GetUserGroupShips(ctx context.Context) ([]uint64, error)
}

func NewWsRepo(db *sql.DB) WsRepo {
	return &WsRepoImpl{db: db}
}

type WsRepoImpl struct {
	db *sql.DB
}

func (w *WsRepoImpl) GetUserFriendShips(ctx context.Context) ([]uint64, error) {
	userId, ok := ctx.Value("user_id").(uint64)
	if !ok {
		utils.GetLogger().Error("类型断言失败")
		return nil, errors.New("类型断言失败")
	}

	query := `select case when user1_id = ? then user2_id else user1_id end as friend_id from friendships where (user1_id = ? or user2_id = ?) and status = 'accepted'`

	rows, err := w.db.QueryContext(ctx, query, userId, userId, userId)
	if err != nil {
		utils.GetLogger().Error("数据库查询失败", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var friendIds []uint64
	for rows.Next() {
		var friendId uint64

		if err := rows.Scan(&friendId); err != nil {
			utils.GetLogger().Error("读取rows失败", zap.Error(err))
			return nil, err
		}

		friendIds = append(friendIds, friendId)
	}
	if err := rows.Err(); err != nil {
		utils.GetLogger().Error("遍历rows失败", zap.Error(err))
		return nil, err
	}

	return friendIds, nil
}

func (w *WsRepoImpl) GetUserGroupShips(ctx context.Context) ([]uint64, error) {
	userId, ok := ctx.Value("user_id").(uint64)
	if !ok {
		utils.GetLogger().Error("类型断言失败")
		return nil, errors.New("类型断言失败")
	}

	query := `select group_id from group_members where user_id = ?`

	rows, err := w.db.QueryContext(ctx, query, userId)
	if err != nil {
		utils.GetLogger().Error("数据库查询失败", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var groupIds []uint64
	for rows.Next() {
		var groupId uint64

		if err := rows.Scan(&groupId); err != nil {
			utils.GetLogger().Error("读取rows失败", zap.Error(err))
			return nil, err
		}

		groupIds = append(groupIds, groupId)
	}
	if err := rows.Err(); err != nil {
		utils.GetLogger().Error("遍历rows失败", zap.Error(err))
		return nil, err
	}

	return groupIds, nil
}
