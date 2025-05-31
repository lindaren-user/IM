package controller

import (
	"IM_BE/Result"
	"IM_BE/dto"
	"IM_BE/service"
	"IM_BE/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (u *UserController) Login(c *gin.Context) {
	var user *dto.UserLoginReqDto

	// TODO：Body io.ReadCloser
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.GetLogger().Error("解析请求体失败", zap.Error(err))
		Result.Error(c, "解析请求体失败")
		return
	}

	token, err := u.service.Login(user.Username, user.Password)
	if err != nil {
		Result.Error(c, "登录失败")
		return
	}

	Result.Success(c, token)
}

func (u *UserController) Logout(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		utils.GetLogger().Error("user_id 不存在")
		Result.Error(c, "登出失败")
		return
	}

	id, ok := userId.(int)
	if !ok {
		utils.GetLogger().Error("断言失败")
		Result.Error(c, "登出失败")
		return
	}

	if err := u.service.Logout(id); err != nil {
		Result.Error(c, "登出失败")
		return
	}

	Result.Success(c, nil)
}
