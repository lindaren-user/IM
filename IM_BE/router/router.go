package router

import (
	"IM_BE/controller"
	"IM_BE/db/mysql"
	"IM_BE/db/redis"
	"IM_BE/middleware"
	"IM_BE/repository"
	"IM_BE/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 允许的前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 允许携带 Cookies
		MaxAge:           12 * time.Hour,
	}))

	repo := repository.NewUserRepo(mysql.Get())
	userService := service.NewUserService(repo, redis.Get())
	userController := controller.NewUserController(userService)

	router.POST("/", userController.Login)

	authRouter := router.Group("/", middleware.JWTAuthMiddleware())
	{
		//authRouter.GET("/ws", userController.StartWS)
		authRouter.GET("/logout", userController.Logout)
	}
}
