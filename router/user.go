package router

import (
	"fmt"
	"gotest/controller/UserController"

	"github.com/gin-gonic/gin"
)

func UserRouter(base string) {
	r := Router.Group("/" + base)

	{
		r.GET("/list", func(ctx *gin.Context) {
			firstname := ctx.DefaultQuery("firstname", "qje")
			fmt.Println(firstname)
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
