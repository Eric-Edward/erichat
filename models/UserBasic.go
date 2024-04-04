package models

import (
	"EriChat/utils"
	"errors"
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
