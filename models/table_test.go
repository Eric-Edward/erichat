package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestChatRoom(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:Tsinghua@/EcChat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	//err = db.AutoMigrate(&ChatRoom{})
	//err2 := db.AutoMigrate(&ChatRoomNumber{})
	err = db.AutoMigrate(&Message{})
	if err != nil {
		fmt.Println("表创建失败")
		return
	}
}
