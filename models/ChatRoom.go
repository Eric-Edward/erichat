package models

import (
	"EriChat/utils"
	"errors"
	"fmt"
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

func CreatePeerChatRoom(channel, u1, u2 string) (string, error) {
	chatRoom := ChatRoom{
		Cid:     uuid.New().String(),
		Channel: channel,
	}
	db := utils.GetMySQLDB()
	tx := db.Begin()
	result := db.Create(&chatRoom)
	if result.RowsAffected != 1 {
		tx.Rollback()
		return "", tx.Error
	}

	//TODO 这里还要添加创建当前聊天的用户已经第一个被邀请的用户 [Finish]
	user1 := ChatRoomMember{
		Cid: chatRoom.Cid,
		Uid: u1,
	}
	user2 := ChatRoomMember{
		Cid: chatRoom.Cid,
		Uid: u2,
	}
	r1 := db.Model(&ChatRoomMember{}).Create(&user1)
	r2 := db.Model(&ChatRoomMember{}).Create(&user2)
	if r1.RowsAffected != 1 || r2.RowsAffected != 1 {
		tx.Rollback()
		fmt.Println("创建用户聊天时失败")
		return "", errors.Join(r1.Error, r2.Error)
	}
	tx.Commit()
	return chatRoom.Cid, nil
}

func GetAllChatRoomByUid(uid string) []ChatRoom {
	var chatRooms []ChatRoom
	var chatRoomNumber ChatRoomMember
	db := utils.GetMySQLDB()
	tx := db.Model(&ChatRoomMember{}).Where("uid=?", uid).Find(&chatRoomNumber)
	if tx.Error != nil {
		return nil
	}
	tx = db.Model(&ChatRoom{}).Where("cid=?", chatRoomNumber.Cid).Find(&chatRooms)
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
