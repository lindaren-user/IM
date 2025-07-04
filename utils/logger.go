package utils

import "go.uber.org/zap"

var Logger *zap.Logger

// TODO：写入log文件

func InitLogger() {
	Logger, _ = zap.NewDevelopment()
}

func GetLogger() *zap.Logger {
	return Logger
}
