package main

import (
	"gotest/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := router.SetupRouter()
	_ = r.Run()
}
