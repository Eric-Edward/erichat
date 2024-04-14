package models

import (
	"EriChat/utils"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ChatRoom struct {
	Cid       string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Channel   string         `gorm:"unique;not null"`
}

type ChatRoomMember struct {
	gorm.Model
	Cid      string   `gorm:"not null;size:191"`
	Uid      string   `gorm:"not null"`
	ChatRoom ChatRoom `gorm:"foreignKey:Cid;references:Cid"`
}

func CreateChatRoom(chatRoomName string, clients []interface{}) (bool, error) {
	chatRoom := ChatRoom{
		Cid:     uuid.New().String(),
		Channel: chatRoomName,
	}
	db := utils.GetMySQLDB()
	tx := db.Begin()
	result := db.Create(&chatRoom)
	if result.Error != nil {
		tx.Rollback()
		return false, tx.Error
	}
	for _, client := range clients {
		result = db.Model(&ChatRoomMember{}).Create(&ChatRoomMember{
			Cid: chatRoom.Cid,
			Uid: client.(string),
		})
		if result.Error != nil {
			tx.Rollback()
			return false, tx.Error
		}
	}
	tx.Commit()
	return true, nil
}

func GetAllChatRoomByUid(uid string) []ChatRoom {
	var chatRooms []ChatRoom
	var chatRoomCid []string
	db := utils.GetMySQLDB()
	tx := db.Model(&ChatRoomMember{}).Select("cid").Where("uid=?", uid).Find(&chatRoomCid)
	if tx.Error != nil {
		return nil
	}
	tx = db.Model(&ChatRoom{}).Where("cid in ?", chatRoomCid).Find(&chatRooms)
	if tx.Error != nil {
		return nil
	}

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

func IsChatRoomMember(cid string, uid string) (bool, error) {
	var chatRoomMember ChatRoomMember
	db := utils.GetMySQLDB()
	result := db.Model(&ChatRoomMember{}).Where("cid=? and uid=?", cid, uid).First(&chatRoomMember)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, nil
}
