package models

import (
	"EriChat/utils"
	"fmt"
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

type RelationShipApply struct {
	gorm.Model
	Apply   string `gorm:"not null"`
	Applied string `gorm:"not null"`
	Group   string `gorm:"not null"`
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

func ApplyRelationShip(uid, fid, group string) (bool, error) {
	var apply = RelationShipApply{Apply: uid, Applied: fid, Group: group}
	db := utils.GetMySQLDB()
	tx := db.Begin()
	result := tx.Model(&RelationShipApply{}).Create(&apply)
	if result.Error != nil {
		fmt.Println("申请失败")
		tx.Rollback()
		return false, tx.Error
	}
	tx.Commit()
	return true, nil
}

func GetRelationShipApplyByUid(applied string) ([]RelationShipApply, error) {
	var apples []RelationShipApply
	db := utils.GetMySQLDB()
	result := db.Model(&RelationShipApply{}).Where("applied=?", applied).Find(&apples)
	if result.Error != nil {
		fmt.Println("查询朋友申请失败")
		return nil, result.Error
	}
	return apples, nil
}

func HandleRelationShipApply(apply, applied, groupApplied string) (bool, error) {
	db := utils.GetMySQLDB()

	var relation = RelationShipApply{}
	db.Model(&RelationShipApply{}).Where("apply=? and applied=?", apply, applied).First(&relation)
	var relationshipApply = RelationShip{
		Uid:   apply,
		Fid:   applied,
		Group: relation.Group,
	}
	var relationshipApplied = RelationShip{
		Uid:   applied,
		Fid:   apply,
		Group: groupApplied,
	}
	tx := db.Begin()
	result := tx.Model(&RelationShip{}).Create(&relationshipApply)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("添加关系失败")
		return false, tx.Error
	}
	result = tx.Model(&RelationShip{}).Create(&relationshipApplied)
	if result.Error != nil {
		tx.Rollback()
		fmt.Println("添加关系失败")
		return false, tx.Error
	}
	_, err := DropRelationShipApply(apply, applied)
	if err != nil {
		tx.Rollback()
		fmt.Println("删除请求关系失败")
		return false, err
	}
	_, err = DropRelationShipApply(applied, apply)
	if err != nil {
		tx.Rollback()
		fmt.Println("删除请求关系失败")
		return false, err
	}
	tx.Commit()
	return true, nil
}

func DropRelationShipApply(apply, applied string) (bool, error) {

	db := utils.GetMySQLDB()
	tx := db.Begin()
	result := tx.Model(&RelationShipApply{}).Where("apply=? and applied=?", apply, applied).Delete(&RelationShipApply{Apply: apply, Applied: applied})
	if result.Error != nil {
		tx.Rollback()
		return false, tx.Error
	}
	result = tx.Model(&RelationShipApply{}).Where("apply=? and applied=?", applied, apply).Delete(&RelationShipApply{Apply: applied, Applied: apply})
	if result.Error != nil {
		tx.Rollback()
		return false, tx.Error
	}
	tx.Commit()
	return true, nil
}
