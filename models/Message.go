package models

import (
	"EriChat/utils"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SendBy  string `gorm:"not null"`
	Content string
	Url     string
	Size    int64
}

// TODO 每一次的getmessage也应该先从redis中获取，如果redis中不存在的话，然后再去mysql中读取

func GetAllMessage(target string) ([]*Message, error) {
	var messages []*Message
	tableName := "messages" + target
	db := utils.GetMySQLDB()
	tx := db.Table(tableName).Model(&Message{}).Find(&messages)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return messages, nil
}
