package UserService

import (
	"EcChat/models/User"
	"EcChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// Register
// @Tags 注册功能
// @Description 处理注册功能
// @Success 200 {string} json:{code:message}
// @router /register [post]
func Register(c *gin.Context) {
	engine := utils.GetMySQLDB()
	user := User.UserBasic{}
	_ = engine.AutoMigrate(&User.UserBasic{})
	err := c.ShouldBind(&user)
	if err != nil {
		fmt.Println("数据转换失败！")
	}
	user.ID = uuid.New().String()
	user.PassWord, err = utils.EncodeInfo(user.PassWord)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "注册失败",
			"state":   "failed",
		})
		return
	}
	engine.Create(&user)
	c.JSON(http.StatusOK, gin.H{
		"message":  "注册成功",
		"uuid":     user.ID,
		"username": user.UserName,
		"state":    "ok",
	})
	return
}

// UsernameIsRegistered
// @Tags 注册功能
// @Description 查看当前用户名是否已经被注册
// @Success 200 {string} json:{code:message}
// @router /register [get]
func UsernameIsRegistered(c *gin.Context) {
	username := c.Query("username")
	engine := utils.GetMySQLDB()
	var user User.UserBasic
	find := engine.Where("user_name=?", username).Find(&user)
	if find.RowsAffected != 0 {
		c.JSON(http.StatusOK, gin.H{
			"state": "failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"state": "ok",
	})
}
