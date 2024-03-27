package User

import (
	"EcChat/utils"
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

// BeforeInsert
// 这个函数在这里已经不需要了，因为xorm可以自己识别CreateAt和UpdateAt
///
/*func (receiver *UserBasic) BeforeInsert() {
	//receiver.CreatedAt = time.Now()
	//receiver.UpdateAt = time.Now()
	fmt.Println("这里执行了", receiver.CreatedAt)
}*/

/*func (receiver *UserBasic) BeforeUpdate() {
	//receiver.UpdateAt = time.Now()
}
*/

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
