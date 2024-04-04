package utils

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

var mySqlDB *gorm.DB
var redisDb *redis.Client

func InitConfig() {
	viper.SetConfigName("app")
	viper.SetConfigFile("config/config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("配置文件读取失败!")
	}

	mySqlDB = getMySQLConnection()
	redisDb = getRedisConnection()
}

func getMySQLConnection() *gorm.DB {
	var s strings.Builder
	s.WriteString(viper.GetString("mysql.user"))
	s.WriteString(":")
	s.WriteString(viper.GetString("mysql.passwd"))
	s.WriteString("@/")
	s.WriteString(viper.GetString("mysql.database"))
	s.WriteString("?charset=utf8mb4&parseTime=True&loc=Local")
	engine, err := gorm.Open(mysql.Open(s.String()), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败！")
	}
	db, _ := engine.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	return engine
}

func GetMySQLDB() *gorm.DB {
	return mySqlDB
}

func getRedisConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "Tsinghua",
		DB:       0,
	})
}

func GetRedis() *redis.Client {
	return redisDb
}
