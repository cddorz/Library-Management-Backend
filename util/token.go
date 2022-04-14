package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var UserKey = []byte("VCisHereUser")
var AdminKey = []byte("VCisHereAdmin")

type Claims struct {
	UserID int
	jwt.StandardClaims
}

func GenToken(userID int, key []byte) string {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "127.0.0.1",  // 签名颁发者
			Subject:   "library token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func AuthToken(tokenString string, key []byte) (int, bool) {
	if tokenString == "" {
		return -1, false
	}
	token, claims, err := ParseToken(tokenString, key)
	if err != nil || !token.Valid {
		return -1, false
	}
	return claims.UserID, true
}

func ParseToken(tokenString string, key []byte) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return key, nil
	})
	return token, Claims, err
}