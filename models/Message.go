package models

import (
	"EriChat/utils"
	"gorm.io/gorm"
	"slices"
)

type Message struct {
	gorm.Model
	Type     string
	Target   utils.Cid
	Message  string
	UserName string
	Uid      utils.Uid
}

func GetMessageByCid(target string, end uint) ([]*Message, error) {
	var messages []*Message
	tableName := "messages_" + target
	db := utils.GetMySQLDB()
	tx := db.Table(tableName).Model(&Message{}).Where("id < ?", end).Order("id desc").Limit(100).Find(&messages)
	if tx.Error != nil {
		return nil, tx.Error
	}
	slices.Reverse(messages)
	return messages, nil
}

func CreateMessageTable(cid string, tx *gorm.DB) error {
	db := utils.GetMySQLDB()
	err := db.Table("messages_" + cid).AutoMigrate(&Message{})
	return err
}
