package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func EncodeInfo(data string) (string, error) {
	generateFromPassword, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("产生密文失败！")
		return "", err
	}
	return string(generateFromPassword), nil
}

func ComparePassword(hashPassWord, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassWord), []byte(password))
	if err != nil {
		return false
	}
	return true
}
