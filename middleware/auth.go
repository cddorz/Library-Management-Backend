package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"lms/util"
	"net/http"
)

type Claims struct {
	UserID int
	jwt.StandardClaims
}

func UserAuth() gin.HandlerFunc {
	return func (c *gin.Context) {
		if userID, ok := auth(c, util.UserKey); ok {
			c.Set("userID", userID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		}
	}
}

func AdminAuth() gin.HandlerFunc {
	return func (c *gin.Context) {
		if userID, ok := auth(c, util.AdminKey); ok {
			c.Set("userID", userID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
		}
	}
}

func auth(c *gin.Context, key []byte) (userID int, ok bool) {
	tokenString := c.PostForm("token")
	userID, ok = util.AuthToken(tokenString, key)
	return
}

