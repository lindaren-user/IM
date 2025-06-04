package utils

import (
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func InitViper() {
	if err := gotenv.Load(); err != nil {
		log.Println("加载 .env 文件失败")
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		panic("读取配置文件失败" + err.Error())
	}
}
