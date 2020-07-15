package router

import (
	"context"
	"gotest/dbs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gotest/controller/Jwt"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func SetupRouter() {
	Router = gin.Default()

	// Router.Use((func() gin.HandlerFunc {
	// 	return func(ctx *gin.Context) {
	// 		fmt.Println(1)
	// 		ctx.Abort()
	// 		fmt.Println(2)
	// 	}
	// }()), (func() gin.HandlerFunc {
	// 	return func(ctx *gin.Context) {
	// 		fmt.Println(3)
	// 		ctx.Next()
	// 		fmt.Println(4)
	// 	}
	// })())

	Router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "index")
	})

	Router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	Router.GET("/test", Jwt.UseJwt)
	Router.GET("/parse", Jwt.ParseJwt)

	UserRouter("user")
	ManagerRouter("manager")

	Router.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "404")
	})

	server := &http.Server{
		Addr:           ":8080",
		Handler:        Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := dbs.InitEnvironment("pro")
	if err != nil {
		log.Println(err)
	}

	// 怎么关怎么开
	defer dbs.Close()

	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-ch
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(cxt)
	if err != nil {
		log.Println("err", err)
	}

}
