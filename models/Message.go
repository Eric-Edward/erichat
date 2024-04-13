package models

import (
	"EriChat/utils"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

const (
	PIC   string = "picture"
	TEXT  string = "text"
	VIDEO string = "video"
)

type Message struct {
	gorm.Model
	SendBy  string `gorm:"not null"`
	Content string
	Url     string
	Size    int64
}

// TODO 每一次的getmessage也应该先从redis中获取，如果redis中不存在的话，然后再去mysql中读取

func AddMessageToTable(room *utils.ChatRoom) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	redis := utils.GetRedis()
	var ctx = context.Background()
	result, err := redis.Get(ctx, room.Cid).Result()
	if err != nil {
		panic(err)
	}
	var messages []utils.WsMessage
	err = json.Unmarshal([]byte(result), &messages)
	if err != nil {
		panic(err)
	}

	var chatRoomName string
	db := utils.GetMySQLDB()
	tx := db.Model(&ChatRoom{}).Select("channel").First(&chatRoomName)
	if tx.Error != nil {
		panic(tx.Error)
	}
	err = db.Table("message_" + chatRoomName).AutoMigrate(&Message{})
	if err != nil {
		panic(err)
	}
}
