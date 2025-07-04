package main

import (
	"IM_BE/db/mysql"
	"IM_BE/db/redis"
	"IM_BE/router"
	"IM_BE/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	utils.InitLogger()
	defer utils.GetLogger().Sync()

	utils.InitViper()

	mysql.Init()
	defer mysql.Close()

	redis.Init()

	utils.InitJWTKey()

	r := gin.Default()

	router.Init(r)

	host := viper.GetString("server.host")
	port := viper.GetString("server.port")

	r.Run(host + ":" + port)
}
