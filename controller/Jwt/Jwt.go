package Jwt

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("小秘密")

func UseJwt(ctx *gin.Context) {
	username := ctx.Query("username")
	token, err := GenToken(username)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	fmt.Println(token)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "success",
		"token": token,
	})
}

func GenToken(username string) (string, error) {
	c := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString(MySecret)
}

func ParseJwt(ctx *gin.Context) {
	tokenStr := ctx.Query("token")
	cla, err := ParseToken(tokenStr)
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":      "success",
		"username": cla.Username,
	})
}

func ParseToken(tokenString string) (mycla *MyClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)
	if token.Valid {
		fmt.Println("合法")
	} else {
		fmt.Println("不合法")
	}
	mycla, ok := token.Claims.(*MyClaims)
	fmt.Println(ok)
	fmt.Println(mycla.Username)
	if ok && token.Valid {
		return
	}
	fmt.Printf("%v", mycla)
	fmt.Println(mycla.StandardClaims.ExpiresAt)
	return nil, errors.New("token 不合法")
}
