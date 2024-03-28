package models

import (
	"gorm.io/gorm"
	"time"
)

type ChatRoom struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	Channel        string         `gorm:"unique;not null"`
	ChatRoomNumber ChatRoomNumber `gorm:"foreignKey:ID"`
}

type ChatRoomNumber struct {
	ID  uint `gorm:"primaryKey"`
	Uid string
}
