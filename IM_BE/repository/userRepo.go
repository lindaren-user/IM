package repository

import (
	"IM_BE/utils"
	"database/sql"
	"go.uber.org/zap"
)

type UserRepo interface {
	AuthUser(username string, password string) (int, error)
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepoImpl{db: db}
}

type userRepoImpl struct {
	db *sql.DB
}

func (u *userRepoImpl) AuthUser(username string, password string) (int, error) {
	var id int

	query := `select id from user where username = ? and password = ?`
	if err := u.db.QueryRow(query, username, password).Scan(&id); err != nil {
		utils.GetLogger().Error("查询失败", zap.Error(err))
		return 0, err
	}

	return id, nil
}
