package models

import (
	"EcChat/utils"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ChatRoom struct {
	Cid            string `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	Channel        string         `gorm:"unique;not null"`
	ChatRoomMember ChatRoomMember `gorm:"foreignKey:Cid"`
}

type ChatRoomMember struct {
	gorm.Model
	Cid string `gorm:"not null;size:191"`
	Uid string `gorm:"not null"`
}

func CreateChatRoom(channel string) error {
	chatRoom := ChatRoom{
		Cid:     uuid.New().String(),
		Channel: channel,
	}
	db := utils.GetMySQLDB()
	tx := db.Create(&chatRoom)
	if tx.RowsAffected != 1 {
		return tx.Error
	}

	//TODO 这里还要添加创建当前聊天的用户已经第一个被邀请的用户
	return nil
}

func GetAllChatRoomByUid(uid string) []ChatRoom {
	var chatRooms []ChatRoom
	db := utils.GetMySQLDB()
	db.Model(&ChatRoomMember{}).Select("Cid").Where("id=?", uid).Find(&chatRooms)
	return chatRooms
}

func GetChatRoomByCid(cid string) (ChatRoom, error) {
	var chatRoom ChatRoom
	db := utils.GetMySQLDB()
	tx := db.Where("cid=?", cid).First(&chatRoom)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return ChatRoom{}, tx.Error
	}
	return chatRoom, nil
}
