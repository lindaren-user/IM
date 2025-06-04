package controller

import (
	"IM_BE/Result"
	"IM_BE/db/redis"
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

// TODO
// 如果不考虑持久化信息，那么实际上只需要做到”转发消息“
// 但是就是要做到持久化，那么需要对消息进行解析

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

	//idAny, exist := c.Get("user_id")
	//if !exist {
	//	utils.GetLogger().Error("id 获取失败", zap.Error(err))
	//	return
	//}
	//
	//id, ok := idAny.(uint64)
	//if !ok {
	//	utils.GetLogger().Error("id 类型错误，断言失败")
	//	return
	//}

	client := ws.NewClient(id, conn)

	w.service.AddClient(client)

	//// TODO:放在后面的妙处,但是这里执行了Close(),导致关闭连接，致使读取消息失败
	// 怎么处理？？？？
	//defer func() {
	//	if err := conn.Close(); err != nil {
	//		utils.GetLogger().Error("ws 断开失败", zap.Error(err))
	//	}
	//	utils.GetLogger().Info("ws 断开成功")
	//
	//	manager.Unregister <- client
	//}()
}
