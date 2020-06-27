package router

import (
	"gotest/controller/UserController"

	"github.com/gin-gonic/gin"
)

func UserRouter(base string) {
	r := Router.Group("/" + base)

	{
		r.GET("/list", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"msg": "ok",
				"data": []string{
					"list1",
					"list2",
				},
			})
		})

		r.POST("/register", UserController.Register)
	}
}
