package main

import (
	"YN/config"
	"YN/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	config.InitConfig() // 初始化配置模块
	log.InitLog()       //初始化日志模块

	r := gin.Default()
	r.POST("/api/SmartWarehous/Update_task_state", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			log.WebLogger.Errorf("Update_task_state ShouldBindJSON err: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "err: Invalid request format",
			})
			return
		}
		log.WebLogger.Infof("%v        \n%+v", c.Request.URL.Path, req)
		log.WebLogger.Infof("")

		c.JSON(http.StatusOK, gin.H{
			"msg": "successful",
		})
	})
	r.Run(":9522")
}
