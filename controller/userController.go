package controller

import (
	"IM_BE/Result"
	"IM_BE/dto"
	"IM_BE/service"
	"IM_BE/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

	expiration := viper.GetInt("token.expiration")

	ctx := c.Request.Context()
	token, err := u.service.Login(ctx, user.Username, user.Password, expiration)
	if err != nil {
		Result.Error(c, "登录失败")
		return
	}

	name := viper.GetString("cookie.name")
	path := viper.GetString("cookie.path")
	domain := viper.GetString("cookie.domain")
	secure := viper.GetBool("cookie.secure")
	httpOnly := viper.GetBool("cookie.httpOnly")

	// cookie 的 maxAge 单位是 s
	c.SetCookie(name, token, expiration*3600, path, domain, secure, httpOnly)

	Result.Success(c, nil)
}

func (u *UserController) Logout(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		utils.GetLogger().Error("user_id 不存在")
		Result.Error(c, "登出失败")
		return
	}

	id, ok := userId.(uint64)
	if !ok {
		utils.GetLogger().Error("断言失败")
		Result.Error(c, "登出失败")
		return
	}

	ctx := c.Request.Context()
	if err := u.service.Logout(ctx, id); err != nil {
		Result.Error(c, "登出失败")
		return
	}

	Result.Success(c, nil)
}

// TODO:可以实现一个取消搜索的功能？？？？？？？？？？/
func (u *UserController) Search(c *gin.Context) {
	getType := c.Query("type")
	keyword := c.Query("keyword")

	if getType != "0" && getType != "1" {
		Result.Error(c, "参数错误")
		return
	}

	if len(keyword) == 0 {
		Result.Error(c, "参数缺少")
		return
	}

	// TODO:前端取消请求，但是没有无法中断（手动ctrl+c可以）
	ctx := c.Request.Context()

	users, err := u.service.Search(ctx, getType, keyword)
	if err != nil {
		Result.Error(c, "搜索失败")
		return
	}

	Result.Success(c, users)
}

func (u *UserController) GetAllFriends(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		utils.GetLogger().Error("user_id 不存在")
		Result.Error(c, "获取失败")
		return
	}

	id, ok := userId.(uint64)
	if !ok {
		utils.GetLogger().Error("断言失败")
		Result.Error(c, "类型断言失败")
		return
	}

	friends, err := u.service.GetAllFriends(c.Request.Context(), id)
	if err != nil {
		Result.Error(c, "获取好友列表失败")
		return
	}

	Result.Success(c, friends)
}
