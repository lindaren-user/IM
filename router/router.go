package router

import (
	"IM_BE/controller"
	"IM_BE/db/mysql"
	"IM_BE/db/redis"
	"IM_BE/middleware"
	"IM_BE/repository"
	"IM_BE/service"
	"IM_BE/ws"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	redisRepo := repository.NewRedisRepo(redis.Get())

	userRepo := repository.NewUserRepo(mysql.Get())
	userService := service.NewUserService(userRepo, redisRepo)
	userController := controller.NewUserController(userService)

	wsRepo := repository.NewWsRepo(mysql.Get())
	wsService := service.NewWsService(wsRepo)
	wsController := controller.NewWsController(wsService)

	ws.InitWsManager()

	router.POST("/user", userController.Login)
	router.GET("/ws", wsController.Run)

	// TODO:中间件是否一定要 c.Next()
	authRouter := router.Group("/user", middleware.JWTAuthMiddleware())
	{
		authRouter.GET("/getAllFriends", userController.GetAllFriends)
		authRouter.GET("/search", userController.Search)
		authRouter.GET("/logout", userController.Logout)
	}
}
