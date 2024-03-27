package utils

import (
	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"time"
)

type UserAuth struct {
	Uid string
	jwt.RegisteredClaims
}

var secretKey = "abcdweHEWFD234Sdfewad"

func init() {
	tick := time.Tick(time.Hour)
	go func() {
		for {
			<-tick
			secret, err := GenerateSecretKey()
			if err != nil {
				return
			}
			secretKey = secret
		}
	}()
}

func GenerateSecretKey() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLength := big.NewInt(int64(len(charset)))
	result := make([]byte, 64)
	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}
	return string(result), nil
}

func GenerateJWT(uid string, t time.Time) (string, error) {
	claims := UserAuth{
		uid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.SignedString([]byte(secretKey))

	return signedString, err
}

func ParseJWT(tokenString string) (*UserAuth, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserAuth{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserAuth); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func IsExpiredJWT(auth *UserAuth) bool {
	expiredTime := auth.ExpiresAt.Time
	if time.Now().After(expiredTime) {
		return true
	} else {
		return false
	}
}
