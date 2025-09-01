package router

import (
	"github.com/gin-gonic/gin"
)

func Url(router *gin.Engine) {
	router.POST("/rbrainrobot/post_field_function", PostFieldFunction) // 1.电梯状态回传、2.独占区配置

	router.POST("/test/read_elevator_status", ReadElevatorStatus) // 读取电梯输入输出的信息
	router.POST("/test/write_to_elevator", WriteToElevator)       // 写入电梯
	router.POST("/test/one_click_reset", OneClickReset)           // 清空电梯任务、清空电梯输入

	router.GET("/api/config", ReadConfigFront)
	router.POST("/api/save_config", SaveConfigFront) // 清空电梯任务、清空电梯输入
}
