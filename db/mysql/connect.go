package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// 单例模式
var db *sql.DB

func Init() {
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	dbname := viper.GetString("mysql.dbname")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", user, password, host, port, dbname)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("数据库打开失败")
	}

	if err = db.Ping(); err != nil {
		panic("数据库无法连接")
	}
}

func Close() {
	if err := db.Close(); err != nil {
		panic("数据库关闭失败")
	}
}

func Get() *sql.DB {
	return db
}
