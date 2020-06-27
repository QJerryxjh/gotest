package UserController

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	params, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		ctx.JSON(200, gin.H{
			"code": "400",
			"msg":  "接受参数错误",
		})
		return
	}

	var data struct {
		Username string `json:"username"`
		Pwd      string `json: "pwd"`
		Email    string `json: "email"`
	}

	err = json.Unmarshal(params, &data)

	if err != nil {
		ctx.JSON(200, gin.H{
			"code": "400",
			"msg":  "接受参数错误2",
		})
		return
	}

	if data.Username == "" {
		ctx.JSON(200, gin.H{
			"code": "400",
			"msg":  "用户名不能为空",
		})
		return
	}

	if data.Pwd == "" {
		ctx.JSON(200, gin.H{
			"code": "400",
			"msg":  "密码不能为空",
		})
		return
	}

	if data.Email == "" {
		ctx.JSON(200, gin.H{
			"code": "400",
			"msg":  "邮箱不能为空",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": "200",
		"msg":  "ok",
	})
}
