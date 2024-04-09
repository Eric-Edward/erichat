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
	Type      string
	Target    Cid
	Message   string
	QuitRooms []string
}

type ChatRoom struct {
	Clients []Uid
	Cid     string
}

type Connection struct {
	Conn         *websocket.Conn
	FromWS       chan WsMessage
	ToWS         chan WsMessage
	CloseReceive chan bool
	CloseSend    chan bool
}

var AllConnections map[Uid]*Connection

var AllChatRooms map[Cid]*ChatRoom

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

	AllConnections = make(map[Uid]*Connection)
	AllChatRooms = make(map[Cid]*ChatRoom)
}

func (conn *Connection) ReceiveEvent() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	for {
		select {
		case msg := <-conn.FromWS:
			room, exist := AllChatRooms[msg.Target]
			if !exist {
				panic("当前聊天室不存在")
			}
			for _, uid := range room.Clients {
				AllConnections[uid].ToWS <- msg
			}
		case <-conn.CloseReceive:
			return
		}
	}
}

func (conn *Connection) SendEvent() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	for {
		select {
		case msg := <-conn.ToWS:
			marshal, _ := json.Marshal(msg)
			err := conn.Conn.WriteMessage(websocket.TextMessage, marshal)
			if err != nil {
				panic(err)
			}
		case <-conn.CloseSend:
			return
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
