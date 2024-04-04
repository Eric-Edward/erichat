package models

import (
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestUserBasic(t *testing.T) {

	engine, err := gorm.Open(mysql.Open("root:Tsinghua@/EcChat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		return
	}
	err = engine.AutoMigrate(&RelationShip{})
	if err != nil {
		return
	}

	//engine.Migrator().DropTable(&UserBasic{})
}
