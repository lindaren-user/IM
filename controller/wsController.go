package controller

import (
	"IM_BE/Result"
	"IM_BE/db/mysql"
	"IM_BE/db/redis"
	"IM_BE/repository"
	"IM_BE/service"
	"IM_BE/utils"
	"IM_BE/ws"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

// TODO：心跳检测

type WsController struct {
	service *service.WsService
}

func NewWsController(service *service.WsService) *WsController {
	return &WsController{service: service}
}

func (w *WsController) Run(c *gin.Context) {
	// TODO:优化校验token代码，怎么写才不冗余？
	token := c.Query("token")

	claims, err := utils.ParseJWT(token)
	if err != nil {
		Result.Error(c, "会话失效") // TODO:会发送？？？
		return
	}

	id := claims.UserID

	tokenKey := fmt.Sprintf("user_%d_token", id)
	redisToken, err := redis.Get().Get(context.Background(), tokenKey).Result()
	if err != nil || redisToken != token {
		Result.Error(c, "会话失效")
		return
	}

	upgrader := websocket.Upgrader{
		// ReadBufferSize:  1024,
		// WriteBufferSize: 1024,
		// 用于检查请求的来源是否允许建立 WebSocket 连接。这里简单返回 true 表示允许所有来源的请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 升级，实际上 http.ResponseWriter 被 WebSocket“劫持”了，也就是：
	// 从正常的 HTTP 响应流程中脱离，变成了 WebSocket 长连接，这时再使用 c.JSON()、Result.Error() 这些基于 HTTP 响应的函数，就会触发该错误。
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.GetLogger().Error("ws 连接失败", zap.Error(err))
		// Result.Error(c, "ws 连接失败") // 不能再调用，会出现 Error #01: http: connection has been hijacked
		return
	}

	messageService := service.NewMessageService(
		repository.NewMessageRepo(mysql.Get()),
		repository.NewRedisRepo(redis.Get()),
	)

	client := ws.NewClient(id, conn, messageService.SaveMessage)

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "user_id", id)
	if err := w.service.AddClient(ctx, client); err != nil {
		return
	}
}
