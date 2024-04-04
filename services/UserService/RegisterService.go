package UserService

import (
	"EriChat/models"
	"EriChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func Register(c *gin.Context) {
	engine := utils.GetMySQLDB()
	user := models.UserBasic{}
	_ = engine.AutoMigrate(&models.UserBasic{})
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
	jwt, _ := utils.GenerateJWT(user.ID, time.Now().Add(time.Minute*20))
	c.JSON(http.StatusOK, gin.H{
		"message":  "注册成功",
		"uuid":     user.ID,
		"username": user.UserName,
		"state":    "ok",
		"token":    jwt,
	})
	return
}

func UsernameIsRegistered(c *gin.Context) {
	username := c.Query("username")
	engine := utils.GetMySQLDB()
	var user models.UserBasic
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
