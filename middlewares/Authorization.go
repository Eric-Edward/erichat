package middlewares

import (
	"EriChat/models"
	"EriChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		splits := strings.Split(token, " ")
		if len(splits) != 2 || strings.ToLower(splits[0]) != "bearer" {
			fmt.Println("Invalid Authorization header format")
			c.JSON(http.StatusOK, gin.H{
				"message": "身份验证失败",
				"code":    utils.FailedBearer,
			})
			c.Abort()
			return
		}
		claims, err := utils.ParseJWT(splits[1])
		if err != nil {
			fmt.Println("Parse JWT failed")
			c.JSON(http.StatusOK, gin.H{
				"message": "登陆超时！",
				"code":    utils.FailedParseJWT,
			})
			c.Abort()
			return
		}
		_, err = models.GetUserByID(claims.Uid)
		isExpiredJWT := utils.IsExpiredJWT(claims)
		if isExpiredJWT || err != nil {
			fmt.Println("User doesn't exist or auth is expired")
			c.JSON(http.StatusOK, gin.H{
				"message": "用户不存在或用户token过期",
				"code":    utils.FailedExpiredJWT,
			})
			c.Abort()
			return
		}
		jwt, err := utils.GenerateJWT(claims.Uid, claims.ExpiresAt.Add(time.Minute*1))
		if err != nil {
			fmt.Println("GenerateJWT make mistake")
			c.JSON(http.StatusOK, gin.H{
				"message": "产生新用户信息时出现错误",
				"code":    utils.FailedGenerateJWT,
			})
			c.Abort()
			return
		}
		c.Header("newToken", jwt)
		c.Set("self", claims.Uid)
		c.Next()
	}
}
