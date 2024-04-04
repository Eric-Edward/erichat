package models

import (
	"EriChat/utils"
	"fmt"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Uid       string `gorm:"index"`
	GroupName string `gorm:"size:64;index"`

	UserBasic UserBasic `gorm:"foreignKey:Uid;references:ID"`
}

func GetAllGroupByUid(uid string) ([]Group, error) {
	var groups []Group
	db := utils.GetMySQLDB()
	tx := db.Model(&Group{}).Where("uid=?", uid).Find(&groups)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return groups, nil
}

func AddGroup(uid, groupName string) (bool, error) {
	var group = Group{Uid: uid, GroupName: groupName}
	db := utils.GetMySQLDB()
	tx := db.Begin()
	result := tx.Model(&Group{}).Create(&group)
	if result.Error != nil || result.RowsAffected != 1 {
		tx.Rollback()
		fmt.Println("添加用户分组失败")
		return false, result.Error
	}
	tx.Commit()
	return true, nil
}

func DeleteGroup(uid, groupName string) (bool, error) {
	var group = Group{Uid: uid, GroupName: groupName}
	db := utils.GetMySQLDB()
	tx := db.Begin()
	result := tx.Delete(&group)
	if result.Error != nil || result.RowsAffected != 1 {
		tx.Rollback()
		fmt.Println("分组删除失败")
		return false, tx.Error
	}
	tx.Commit()
	return true, nil
}
