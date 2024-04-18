package models

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestChatRoom(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:Tsinghua@/EcChat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	//err = db.AutoMigrate(&UserBasic{})
	//err = db.AutoMigrate(&Group{})

	err = db.AutoMigrate(&ChatRoom{})
	//err = db.AutoMigrate(&ChatRoomMember{})
	//err = db.AutoMigrate(&RelationShip{})
	if err != nil {
		fmt.Println("表创建失败")
		return
	}
}

func TestCreateChatRoom(t *testing.T) {
	chatRoom := ChatRoom{
		Cid:     uuid.New().String(),
		Channel: "test",
	}
	db, _ := gorm.Open(mysql.Open("root:Tsinghua@/EcChat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	tx := db.Create(&chatRoom)
	if tx.RowsAffected != 1 {
		fmt.Println("生成失败")
	}
}

func TestCreateChatRoomNumber(t *testing.T) {
	member := ChatRoomMember{
		Cid: "bbe50b65-5eda-4182-a44f-5f",
		Uid: "111",
	}
	db, _ := gorm.Open(mysql.Open("root:Tsinghua@/EcChat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	tx := db.Create(&member)
	if tx.RowsAffected != 1 {
		fmt.Println("插入聊天人员失败！")
	}
}
