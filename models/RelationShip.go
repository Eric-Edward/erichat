package models

import (
	"EriChat/utils"
	"gorm.io/gorm"
)

type RelationShip struct {
	gorm.Model
	Uid    string `gorm:"not null"`
	Fid    string `gorm:"not null"`
	Remark string
	Group  string `gorm:"not null"`
	Groups Group  `gorm:"foreignKey:Group;references:GroupName"`
}

func GetRelationShipByUid(uid string) ([]RelationShip, error) {
	var relationships []RelationShip
	db := utils.GetMySQLDB()
	tx := db.Model(&RelationShip{}).Where("uid=?", uid).Find(&relationships)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return relationships, nil
}

func AddRelationShipByUid(relation RelationShip) (bool, error) {
	var user UserBasic
	db := utils.GetMySQLDB()
	tx := db.Begin()
	result := tx.Model(&UserBasic{}).Where("id=?", relation.Uid).First(&user)
	if result.Error != nil {
		tx.Rollback()
		return false, result.Error
	}
	result = tx.Model(&UserBasic{}).Where("id=?", relation.Fid).First(&user)
	if result.Error != nil {
		tx.Rollback()
		return false, tx.Error
	}
	if relation.Remark == "" {
		relation.Remark = user.UserName
	}
	result = tx.Model(&RelationShip{}).Create(relation)
	if result.Error != nil {
		tx.Rollback()
		return false, tx.Error
	}
	tx.Commit()
	return true, nil
}

func DeleteRelationShipByUid(relation RelationShip) (bool, error) {
	db := utils.GetMySQLDB()
	tx := db.Begin()
	result := tx.Model(&RelationShip{}).Delete(relation)
	if result.Error != nil {
		tx.Rollback()
		return false, tx.Error
	}
	tx.Commit()
	return true, nil
}
