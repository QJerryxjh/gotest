package router

import (
	"context"
	"fmt"
	"gotest/dbs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func SetupRouter() {
	Router = gin.Default()

	Router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "index")
	})

	Router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	Router.GET("/test", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		options := sessions.Options{
			Path:   "/",
			MaxAge: 604800,
		}
		session.Options(options)
		err := session.Save()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(session)
		ctx.JSON(200, gin.H{
			"msg": "success",
		})
	})

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
