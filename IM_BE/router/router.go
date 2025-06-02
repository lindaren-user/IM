package router

import (
	"IM_BE/controller"
	"IM_BE/db/mysql"
	"IM_BE/db/redis"
	"IM_BE/middleware"
	"IM_BE/repository"
	"IM_BE/service"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	userRepo := repository.NewUserRepo(mysql.Get())
	userService := service.NewUserService(userRepo, redis.Get())
	userController := controller.NewUserController(userService)

	wsRepo := repository.NewWsRepo(mysql.Get())
	wsService := service.NewWsService(wsRepo)
	wsController := controller.NewWsController(wsService)

	router.POST("/", userController.Login)

	// TODO:中间件是否一定要 c.Next()
	authRouter := router.Group("/", middleware.JWTAuthMiddleware())
	{
		authRouter.GET("/logout", userController.Logout)
	}

	router.GET("/ws", wsController.Work)
}
