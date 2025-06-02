package Result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// TODO:http响应码默认是200，到底怎么设置
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Result{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, &Result{
		Code:    1,
		Message: msg,
	})
}
