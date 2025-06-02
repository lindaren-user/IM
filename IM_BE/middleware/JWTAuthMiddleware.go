package middleware

import (
	"IM_BE/Result"
	"IM_BE/db/redis"
	"IM_BE/utils"
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			Result.Error(c, "Authorization header missing or invalid")
			// c.Abort() 只是可以使得 index 超出范围，还需要配合 return 退出当前函数
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseJWT(token)
		if err != nil {
			Result.Error(c, "会话失效")
			c.Abort()
			return
		}

		id := claims.UserID

		tokenKey := fmt.Sprintf("user_token_%d", id)
		tokenTmp, err := redis.Get().Get(context.Background(), tokenKey).Result()
		if err != nil || tokenKey != tokenTmp {
			Result.Error(c, "会话失效")
			c.Abort()
			return
		}

		// 把解析到的用户 ID 存入 context 中
		c.Set("user_id", id)

		// c.Next() // 多此一举
	}
}
