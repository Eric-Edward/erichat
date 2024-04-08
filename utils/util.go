package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

type Uid string
type Cid string

type WsMessage struct {
	Target  Cid
	Message string
}

type ChatRoom struct {
	Clients []Uid
	Cid     string
}

type Connection struct {
	Conn   *websocket.Conn
	FromWS chan WsMessage
	ToWS   chan WsMessage
}

var AllConnections map[Uid]*Connection

var AllChatRooms map[Cid]*ChatRoom

var mySqlDB *gorm.DB

var redisDb *redis.Client

var channel *DeliverMessage

func InitConfig() {
	viper.SetConfigName("app")
	viper.SetConfigFile("config/config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("配置文件读取失败!")
	}

	mySqlDB = getMySQLConnection()
	redisDb = getRedisConnection()

	channel = &DeliverMessage{
		Channel: make(chan string),
		Message: make(chan []byte),
	}

	AllConnections = make(map[Uid]*Connection)
	AllChatRooms = make(map[Cid]*ChatRoom)
}

func (conn *Connection) EventLoop() {
	for {
		select {
		case msg := <-conn.ToWS:
			marshal, _ := json.Marshal(msg)
			err := conn.Conn.WriteMessage(websocket.TextMessage, marshal)
			if err != nil {

			}
		case msg := <-conn.FromWS:
			for _, uid := range AllChatRooms[msg.Target].Clients {
				AllConnections[uid].ToWS <- msg
			}
		}

	}
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

func GetChannel() *DeliverMessage {
	return channel
}
