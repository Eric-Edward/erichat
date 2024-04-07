package models

import (
	"EriChat/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserName string `gorm:"unique;required"`
	PassWord string `gorm:"required"`
	Age      int
	Phone    string
	Email    string
}

type UserInfo struct {
	ID       string `gorm:"primarykey"`
	UserName string `gorm:"unique;notnull"`
	Age      int
	Phone    string `validate:"isPhoneNumber"`
	Email    string `validate:"email"`
}

type Client struct {
	ID       string
	UserName string
}

func GetUserByID(id string) (UserBasic, error) {
	var user UserBasic
	db := utils.GetMySQLDB()
	tx := db.Where("id=?", id).First(&user)
	exist := errors.Is(tx.Error, gorm.ErrRecordNotFound)
	if !exist {
		return user, nil
	}
	return user, tx.Error
}

func GetAllClientsByUserName(username string) ([]Client, error) {
	var clients []Client
	db := utils.GetMySQLDB()
	tx := db.Model(&UserBasic{}).Where("user_name like ?", "%"+username+"%").Find(&clients)
	if tx.Error != nil {
		fmt.Println("查找失败")
		return nil, tx.Error
	}
	return clients, nil
}
