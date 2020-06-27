package router

import (
	"gotest/controller/ManagerController"
)

func ManagerRouter(base string) {
	r := Router.Group("/" + base)

	{
		r.GET("/list", ManagerController.GetManagerList)
	}

}
