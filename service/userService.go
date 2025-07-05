package service

import (
	"IM_BE/dto"
	"IM_BE/repository"
	"IM_BE/utils"
	"context"
	"go.uber.org/zap"
)

type UserService struct {
	repo      repository.UserRepo
	redisRepo repository.RedisRepo
}

func NewUserService(repo repository.UserRepo, redisRepo repository.RedisRepo) *UserService {
	return &UserService{repo: repo, redisRepo: redisRepo}
}

func (u *UserService) Login(ctx context.Context, username string, password string, expiration int) (string, error) {
	id, err := u.repo.AuthUser(ctx, username, password)
	if err != nil {
		return "", nil
	}

	token, err := utils.GenerateJWT(id)
	if err != nil {
		utils.GetLogger().Error("获取 token 失败", zap.Error(err))
		return "", nil
	}

	if err := u.redisRepo.SetUserToken(ctx, id, token, expiration); err != nil {
		utils.GetLogger().Error("redis 获取 token 失败", zap.Error(err))
		return "", err
	}

	return token, nil
}

func (u *UserService) Logout(ctx context.Context, id uint64) error {
	if err := u.redisRepo.DelUserToken(ctx, id); err != nil {
		utils.GetLogger().Error("redis 删除 token 失败", zap.Error(err))
		return err
	}

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

func (u *UserService) GetAllFriends(ctx context.Context, id uint64) ([]*dto.FriendDTO, error) {
	friendIds, err := u.repo.GetAllFriends(ctx, id)
	if err != nil {
		return nil, err
	}

	var friends []*dto.FriendDTO
	for _, friendId := range friendIds {
		friend, err := u.repo.GetFriendInfo(ctx, friendId)
		if err != nil {
			return nil, err
		}

		friends = append(friends, friend)
	}

	return friends, nil
}
