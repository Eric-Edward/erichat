package UserService

import (
	"EriChat/models"
	"EriChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"regexp"
)

func CompleteUserInfo(c *gin.Context) {
	var user models.UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		fmt.Println("数据绑定失败")
		c.JSON(http.StatusOK, gin.H{
			"message": "数据绑定失败",
			"state":   "failed",
		})
		return
	}
	db := utils.GetMySQLDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.RegisterValidation("isPhoneNumber", ValidateMyPhone)
	if err != nil {
		return
	}
	validatorErr := validate.Struct(user)
	if validatorErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "输入的数据格式不正确",
			"state":   "failed",
		})
		tx.Rollback()
		return
	}
	updates := tx.Model(models.UserBasic{}).Where("id=?", user.ID).Updates(map[string]interface{}{
		"phone": user.Phone,
		"email": user.Email,
	})
	if updates.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "数据更新失败",
			"state":   "failed",
		})
		tx.Rollback()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据更新成功!",
		"state":   "ok",
	})
	tx.Commit()
}

func GetUserInfo(c *gin.Context) {
	uid := c.Query("uuid")
	db := utils.GetMySQLDB()

	var findUser models.UserInfo
	db.Model(&models.UserBasic{}).Where("id=?", uid).Limit(1).Find(&findUser)
	if findUser.ID != uid {
		c.JSON(http.StatusOK, gin.H{
			"message": "数据库中不能存在当前用户",
			"uuid":    uid,
			"state":   "failed",
		})
		return
	}

	c.JSON(http.StatusOK, findUser)
}

func ValidateMyPhone(level validator.FieldLevel) bool {
	matched, err := regexp.Match("^1[3-9]\\d{9}$", []byte(level.Field().String()))
	if err != nil {
		return false
	}
	return matched
}
