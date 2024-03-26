package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.SetConfigFile("config/config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("配置文件读取失败!")
	}
}

func getMySQLConnection() string {
	var s strings.Builder
	s.WriteString(viper.GetString("mysql.user"))
	s.WriteString(":")
	s.WriteString(viper.GetString("mysql.passwd"))
	s.WriteString("@/")
	s.WriteString(viper.GetString("mysql.database"))
	s.WriteString("?charset=utf8mb4&parseTime=True&loc=Local")
	return s.String()

}

func GetMySQLDB() *gorm.DB {
	engine, err := gorm.Open(mysql.Open(getMySQLConnection()), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败！")
	}
	db, _ := engine.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	return engine
}
