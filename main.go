package main

import (
	"gotest/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router.SetupRouter()
}
