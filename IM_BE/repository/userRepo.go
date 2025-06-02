package repository

import (
	"IM_BE/utils"
	"database/sql"
	"go.uber.org/zap"
)

type UserRepo interface {
	AuthUser(username string, password string) (uint64, error)
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepoImpl{db: db}
}

type userRepoImpl struct {
	db *sql.DB
}

func (u *userRepoImpl) AuthUser(username string, password string) (uint64, error) {
	var id uint64

	query := `select id from users where username = ? and password = ?`
	if err := u.db.QueryRow(query, username, password).Scan(&id); err != nil {
		utils.GetLogger().Error("查询失败", zap.Error(err))
		return 0, err
	}

	return id, nil
}
