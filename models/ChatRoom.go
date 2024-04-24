package models

import (
	"EriChat/utils"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ChatRoom struct {
	Cid        string `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Channel    string         `gorm:"unique;not null"`
	Type       string         `gorm:"not null"`
	Avatar     string
	AvatarType string
}

type ChatRoomMember struct {
	gorm.Model
	Cid      string `gorm:"not null;size:191"`
	Uid      string `gorm:"not null"`
	Record   uint
	ChatRoom ChatRoom `gorm:"foreignKey:Cid;references:Cid"`
}

type ChatRoomAside struct {
	Cid        string `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Channel    string         `gorm:"unique;not null"`
	Type       string         `gorm:"not null"`
	Avatar     string
	AvatarType string
}

func CreateChatRoom(chatRoomName string, clients []interface{}) (bool, error) {
	chatRoom := ChatRoom{
		Cid:     uuid.New().String(),
		Channel: chatRoomName,
		Type:    "group",
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
			Cid:    chatRoom.Cid,
			Uid:    client.(string),
			Record: 0,
		})
		if result.Error != nil {
			tx.Rollback()
			return false, tx.Error
		}
	}
	err := tx.Table("messages_" + chatRoom.Cid).AutoMigrate(&Message{})
	if err != nil {
		tx.Rollback()
		return false, tx.Error
	}
	tx.Commit()
	return true, nil
}

func GetAllChatGroupByUid(uid string) []ChatRoom {
	var chatRooms []ChatRoom
	var chatRoomCid []string
	db := utils.GetMySQLDB()
	tx := db.Model(&ChatRoomMember{}).Select("cid").Where("uid=?", uid).Find(&chatRoomCid)
	if tx.Error != nil {
		return nil
	}
	tx = db.Model(&ChatRoom{}).Where("cid in ? and type = 'group'", chatRoomCid).Find(&chatRooms)
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

func UpdateChatRoom(chatRoom ChatRoom) (bool, error) {
	db := utils.GetMySQLDB()
	result := db.Model(&ChatRoom{}).Where("cid=?", chatRoom.Cid).Update("channel", chatRoom.Channel)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func UploadChatRoomAvatar(chatRoom ChatRoom) (bool, error) {
	db := utils.GetMySQLDB()
	tx := db.Begin()
	result := tx.Model(&ChatRoom{}).Where("cid=?", chatRoom.Cid).Updates(chatRoom)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		return false, errors.Join(tx.Error, errors.New("更新失败"))
	}
	tx.Commit()
	return true, nil
}

func GetChatRoomMessageDivider(uid, cid string) (uint, error) {
	var record uint
	db := utils.GetMySQLDB()
	tx := db.Model(&ChatRoomMember{}).Select("record").Where("uid=? and cid=?", uid, cid).First(&record)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return 0, tx.Error
	}
	return record, nil
}
