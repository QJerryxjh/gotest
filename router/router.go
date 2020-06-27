package router

import (
	"gotest/dbs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func SetupRouter() *gin.Engine {
	Router = gin.Default()

	Router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "index")
	})

	Router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	Router.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "404")
	})

	UserRouter("user")
	ManagerRouter("manager")

	err := dbs.InitEnvironment("pro")
	if err != nil {
		log.Println(err)
	}

	// 怎么关怎么开
	defer dbs.Close()

	return Router
}
