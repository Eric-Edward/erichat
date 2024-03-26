package UserService

import (
	"EcChat/models/User"
	"EcChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler
// @Tags 登陆功能
// @Description 这是一个登陆处理的方法
// @Success 200 {string} json:{code:message}
// @router /login [post]
func LoginHandler(c *gin.Context) {
	var user userLogin
	err := c.ShouldBind(&user)
	if err != nil {
		fmt.Println("数据绑定失败：", err)
	}

	var findUser User.UserBasic
	engine := utils.GetMySQLDB()
	_ = engine.AutoMigrate(&User.UserBasic{})
	engine.Select("pass_word").Where("user_name=?", user.Username).Find(&findUser)

	if !utils.ComparePassword(findUser.PassWord, user.Password) {
		c.JSON(http.StatusOK, gin.H{
			"message": "用户不存在或密码不正确",
			"uuid":    "undefined",
			"state":   "failed",
		})
		return
	}

	findResult := engine.Where("user_name=? AND pass_word=?", user.Username, findUser.PassWord).Find(&findUser)

	if findResult.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "登陆失败，数据库中没有该用户",
			"uuid":    "undefined",
			"state":   "failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "登陆成功",
		"status":   "ok",
		"username": user.Username,
		"uuid":     findUser.ID,
	})
	return
}
