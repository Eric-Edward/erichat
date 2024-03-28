package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestUserBasic(t *testing.T) {

	engine, err := gorm.Open(mysql.Open("root:Tsinghua@tcp(127.0.0.1:3306)/EcChat"), &gorm.Config{})
	if err != nil {
		return
	}
	engine.AutoMigrate(&UserBasic{})

	user := UserBasic{UserName: "Eric", PassWord: "Tsinghua", Age: 23}
	result := engine.Create(&user)
	fmt.Println(user.ID, result.RowsAffected)
	var user2 UserBasic
	engine.First(&user2)

	fmt.Println(user2)

	//engine.Migrator().DropTable(&UserBasic{})
}
