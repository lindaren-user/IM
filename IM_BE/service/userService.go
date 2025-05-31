package service

import (
	"IM_BE/repository"
	"IM_BE/utils"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

type UserService struct {
	repo     repository.UserRepo
	redisCli *redis.Client
}

func NewUserService(repo repository.UserRepo, redisCli *redis.Client) *UserService {
	return &UserService{repo: repo, redisCli: redisCli}
}

func (u *UserService) Login(username string, password string) (string, error) {
	id, err := u.repo.AuthUser(username, password)
	if err != nil {
		return "", nil
	}

	token, err := utils.GenerateJWT(id)
	if err != nil {
		utils.GetLogger().Error("获取 token 失败", zap.Error(err))
		return "", nil
	}

	tokenKey := fmt.Sprintf("user_token_%d", id)
	expiration := viper.GetInt("token.expiration")

	// 设置时间，自动清除 token，不占内存
	u.redisCli.Set(context.Background(), tokenKey, token, time.Duration(expiration))

	return token, nil
}

func (u *UserService) Logout(id int) error {
	tokenKey := fmt.Sprintf("user_token_%d", id)

	u.redisCli.Del(context.Background(), tokenKey)

	return nil
}
