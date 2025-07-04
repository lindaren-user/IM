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
