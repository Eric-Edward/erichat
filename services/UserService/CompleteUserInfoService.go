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
			"code":    utils.FailedBindInfo,
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
			"code":    utils.FailedBindInfo,
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
			"code":    utils.FailedUpdateUserInfo,
		})
		tx.Rollback()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "数据更新成功!",
		"code":    utils.Success,
	})
	tx.Commit()
}

func GetUserInfo(c *gin.Context) {
	uid, _ := c.Get("self")
	cid := c.Query("cid")
	db := utils.GetMySQLDB()
	if cid != "" {
		var fid string
		db.Model(&models.ChatRoomMember{}).Select("uid").Where("cid = ? and uid <> ?", cid, uid).First(&fid)
		uid = fid
	}
	var findUser models.UserInfo
	db.Model(&models.UserBasic{}).Where("id=?", uid).Limit(1).Find(&findUser)
	if findUser.ID != uid {
		c.JSON(http.StatusOK, gin.H{
			"message": "数据库中不能存在当前用户",
			"uuid":    uid,
			"code":    utils.FailedFindUser,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取当前用户信息成功",
		"code":    utils.Success,
		"info":    findUser,
	})
}

func ValidateMyPhone(level validator.FieldLevel) bool {
	matched, err := regexp.Match("^1[3-9]\\d{9}$", []byte(level.Field().String()))
	if err != nil {
		return false
	}
	return matched
}

func CompleteGroupInfo(c *gin.Context) {
	var room models.ChatRoom
	err := c.ShouldBind(&room)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "绑定数据失败",
			"code":    utils.FailedBindInfo,
		})
		return
	}
	ok, err := models.UpdateChatRoom(room)
	if err != nil || !ok {
		c.JSON(http.StatusOK, gin.H{
			"message": "更新群聊数据失败",
			"code":    utils.FailedUpdateGroupInfo,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "更新群聊数据成功",
		"code":    utils.Success,
	})
	return
}

func UploadUserAvatar(c *gin.Context) {
	uid, _ := c.Get("self")
	var avatar models.UserAvatar
	err := c.ShouldBind(&avatar)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "绑定数据失败",
			"code":    utils.FailedBindInfo,
		})
		return
	}
	ok, err := models.UpdateUserAvatar(uid.(string), avatar)
	if err != nil || !ok {
		c.JSON(http.StatusOK, gin.H{
			"message": "更新用户头像失败",
			"code":    utils.FailedUpdateUserInfo,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "成功更新用户头像",
		"code":    utils.Success,
	})
	return
}

func GetUserAvatar(c *gin.Context) {
	fid := c.Query("fid")
	avatar, err := models.GetUserAvatarByID(fid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "获取用户头像失败",
			"code":    utils.FailedLoadUserAvatar,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户头像成功",
		"code":    utils.Success,
		"avatar":  avatar,
	})
	return

}
