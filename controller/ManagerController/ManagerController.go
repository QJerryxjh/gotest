package ManagerController

import (
	"gotest/dbs"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetManagerList(ctx *gin.Context) {
	var dbManager dbs.DbManger

	dbManager, err := dbManager.QueryAllManager()
	if err != nil && err == gorm.ErrRecordNotFound {
		ctx.JSON(200, gin.H{
			"error": map[string]interface{}{
				"code": "01003",
				"msg":  "不存在",
			},
		})
		ctx.Abort()
		return
	} else if err != nil {
		ctx.JSON(200, gin.H{
			"error": map[string]interface{}{
				"code": "500",
				"msg":  err.Error(),
			},
		})
		return
	}
}
