package models

import "gorm.io/gorm"

const (
	PIC   string = "picture"
	TEXT  string = "text"
	VIDEO string = "video"
)

type Message struct {
	gorm.Model
	SendBy    string `gorm:"not null"`
	ReceiveBy string `gorm:"not null"`
	Type      string `gorm:"not null"`
	Content   string
	Url       string
	Size      int64
}
