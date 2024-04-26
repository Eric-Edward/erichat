package models

import (
	"EriChat/utils"
	"gorm.io/gorm"
	"slices"
	"time"
)

type Message struct {
	ID        uint `gorm:"primarykey;autoIncrement:false" `
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Type      string
	Target    utils.Cid
	Message   string
	UserName  string
	Uid       utils.Uid
}

func GetMessageByCid(target string, end uint, limit int) ([]*Message, error) {
	var messages []*Message
	tableName := "messages_" + target
	db := utils.GetMySQLDB()
	tx := db.Table(tableName).Model(&Message{}).Where("id < ?", end).Order("id desc").Limit(limit).Find(&messages)
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
