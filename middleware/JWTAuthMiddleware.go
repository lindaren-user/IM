package middleware

import (
	"IM_BE/Result"
	"IM_BE/db/redis"
	"IM_BE/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := viper.GetString("cookie.name")

		token, err := c.Cookie(name)
		if err != nil {
			Result.Error(c, "获取 cookie 失败")
			c.Abort()
			return
		}

		claims, err := utils.ParseJWT(token)
		if err != nil {
			Result.Error(c, "会话失效")
			c.Abort()
			return
		}

		id := claims.UserID

		tokenKey := fmt.Sprintf("user_token_%d", id)
		redisToken, err := redis.Get().Get(context.Background(), tokenKey).Result()
		if err != nil || token != redisToken {
			utils.GetLogger().Error("会话失效", zap.Error(err))
			Result.Error(c, "会话失效")
			c.Abort()
			return
		}

		// 把解析到的用户 ID 存入 context 中
		c.Set("user_id", id)

		// c.Next() // 多此一举
	}
}
