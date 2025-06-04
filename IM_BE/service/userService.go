package service

import (
	"IM_BE/dto"
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

func (u *UserService) Login(ctx context.Context, username string, password string) (string, error) {
	id, err := u.repo.AuthUser(ctx, username, password)
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
	u.redisCli.Set(ctx, tokenKey, token, time.Duration(expiration)*time.Hour)

	fmt.Println(tokenKey, token)

	return token, nil
}

func (u *UserService) Logout(ctx context.Context, id int) error {
	tokenKey := fmt.Sprintf("user_token_%d", id)

	u.redisCli.Del(ctx, tokenKey)

	return nil
}

func (u *UserService) Search(ctx context.Context, getType string, keyword string) (users []*dto.UserDTO, err error) {
	if getType == "0" {
		if users, err = u.repo.GetUserByUsername(ctx, keyword); err != nil {
			return nil, err
		}
	} else {
		if users, err = u.repo.GetUserByNickname(ctx, keyword); err != nil {
			return nil, err
		}
	}

	return users, nil
}
