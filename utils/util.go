package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"sync"
)

type Uid string
type Cid string

type WsMessage struct {
	Type      string
	Target    Cid
	Message   string
	QuitRooms []string `gorm:"null"`
	UserName  string
	Uid       Uid
	Time      string
}

type ChatRoom struct {
	Clients map[Uid]Uid
	Cid     string
	Message []WsMessage
}

type Connection struct {
	Conn         *websocket.Conn
	FromWS       chan WsMessage
	ToWS         chan WsMessage
	CloseReceive chan bool
	CloseSend    chan bool
}

var AllConnections sync.Map

var AllChatRooms sync.Map

var mySqlDB *gorm.DB

var redisDB *redis.Client

func InitConfig() {
	viper.SetConfigName("app")
	viper.SetConfigFile("config/config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("配置文件读取失败!")
	}

	mySqlDB = getMySQLConnection()
	redisDB = getRedisConnection()

	AllConnections = sync.Map{}
	AllChatRooms = sync.Map{}
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
			room, exist := AllChatRooms.Load(msg.Target)
			if !exist {
				panic("当前聊天室不存在")
			}
			for _, uid := range room.(*ChatRoom).Clients {
				connection, _ := AllConnections.Load(uid)
				connection.(*Connection).ToWS <- msg
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
	var conn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "Tsinghua",
		DB:       0,
	})
	var ctx = context.Background()
	conn.ConfigSet(ctx, "maxmemory", "100mb")
	conn.ConfigSet(ctx, "maxmemory-policy", "allkeys-lfu")

	return conn
}

func GetRedis() *redis.Client {
	return redisDB
}
