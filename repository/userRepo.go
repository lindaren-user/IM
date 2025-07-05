package repository

import (
	"IM_BE/dto"
	"IM_BE/utils"
	"context"
	"database/sql"
	"go.uber.org/zap"
)

type UserRepo interface {
	AuthUser(ctx context.Context, username string, password string) (uint64, error)

	GetUserByUsername(ctx context.Context, keyword string) ([]*dto.UserDTO, error)

	GetUserByNickname(ctx context.Context, keyword string) ([]*dto.UserDTO, error)

	GetAllFriends(ctx context.Context, userId uint64) ([]uint64, error)

	GetFriendInfo(ctx context.Context, friendId uint64) (*dto.FriendDTO, error)
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepoImpl{db: db}
}

type userRepoImpl struct {
	db *sql.DB
}

func (u *userRepoImpl) AuthUser(ctx context.Context, username string, password string) (uint64, error) {
	var id uint64

	query := `select id from users where username = ? and password = ?`
	if err := u.db.QueryRowContext(ctx, query, username, password).Scan(&id); err != nil {
		utils.GetLogger().Error("查询失败", zap.Error(err))
		return 0, err
	}

	return id, nil
}

func (u *userRepoImpl) GetUserByUsername(ctx context.Context, keyword string) ([]*dto.UserDTO, error) {
	query := "select id, username, nickname from users where username like ?"

	// 如果查询被中断，返回 context canceled 错误
	rows, err := u.db.QueryContext(ctx, query, keyword+"%")
	if err != nil {
		utils.GetLogger().Error("查询失败", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var users []*dto.UserDTO
	for rows.Next() {
		var user = &dto.UserDTO{}
		if err = rows.Scan(&user.Id, &user.Username, &user.Nickname); err != nil {
			utils.GetLogger().Error("读取失败", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		utils.GetLogger().Error("遍历行失败", zap.Error(err))
		return nil, err
	}

	return users, nil
}

func (u *userRepoImpl) GetUserByNickname(ctx context.Context, keyword string) ([]*dto.UserDTO, error) {
	query := "select id, username, nickname from users where nickname like ?"

	rows, err := u.db.QueryContext(ctx, query, keyword+"%")
	if err != nil {
		utils.GetLogger().Error("查询失败", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var users []*dto.UserDTO
	for rows.Next() {
		var user = &dto.UserDTO{}
		if err = rows.Scan(&user.Id, &user.Username, &user.Nickname); err != nil {
			utils.GetLogger().Error("读取失败", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		utils.GetLogger().Error("遍历行失败", zap.Error(err))
		return nil, err
	}

	return users, nil
}

// TODO:case when then ？？？
func (u *userRepoImpl) GetAllFriends(ctx context.Context, userId uint64) ([]uint64, error) {
	query := `select case when user1_id = ? then user2_id else user1_id end as friend_id from friendships where (user1_id = ? or user2_id = ?) and status = 'accepted'`

	rows, err := u.db.QueryContext(ctx, query, userId, userId, userId)
	if err != nil {
		utils.GetLogger().Error("数据库查询失败", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var friendIds []uint64

	for rows.Next() {
		var friendId uint64
		if err := rows.Scan(&friendId); err != nil {
			utils.GetLogger().Error("读取 rows 失败", zap.Error(err))
			return nil, err
		}

		friendIds = append(friendIds, friendId)
	}
	if err := rows.Err(); err != nil {
		utils.GetLogger().Error("遍历 rows 失败", zap.Error(err))
		return nil, err
	}

	return friendIds, nil
}

func (u *userRepoImpl) GetFriendInfo(ctx context.Context, friendId uint64) (*dto.FriendDTO, error) {
	query := `select id, nickname, avatar from users where id = ?`

	var friend dto.FriendDTO
	var avatar sql.NullString
	if err := u.db.QueryRowContext(ctx, query, friendId).Scan(&friend.Id, &friend.Nickname, &avatar); err != nil {
		utils.GetLogger().Error("数据库查询出错", zap.Error(err))
		return nil, err
	}

	// 注意空字段处理
	friend.Avatar = avatar.String

	return &friend, nil
}
