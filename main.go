package main

import (
	"YN/config"
	"YN/log"
	"YN/router"
	"YN/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig() // 初始化配置模块
	log.InitLog()       //初始化日志模块

	go service.StartEquipmentMonitor() // 读取电梯状态，做逻辑处理
	go service.RasterExclusiveAreaProcess()

	r := gin.Default()
	r.Use(router.Core) // 跨域问题【必须放在url前执行，否则不生效】
	router.Url(r)

	if err := r.Run(":9521"); err != nil {
		log.Logger.Errorf("web start failed, err: %v", err)
	}
}
