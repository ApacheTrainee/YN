package router

import (
	"github.com/gin-gonic/gin"
)

func Url(router *gin.Engine) {
	router.POST("/rbrainrobot/post_field_function", PostFieldFunction) // 1.电梯状态回传、2.独占区配置

	router.POST("/test/read_elevator_status", ReadElevatorStatus)           // 读取电梯输入输出的信息
	router.POST("/test/write_to_elevator", WriteToElevator)                 // 写入电梯
	router.POST("/test/one_click_reset", OneClickReset)                     // 清空电梯任务、清空电梯输入
	router.POST("/test/simulate_write_elevator_read", SimulateElevatorRead) // 因为电梯modbus协议的的input是只读的，只能手动修改，太麻烦，这里模拟，好自动化测试
	router.POST("/test/get_elevator_task", GetElevatorTask)                 // 查询电梯任务状态，用于自动化测试

	router.GET("/api/config", ReadConfigFront)
	router.POST("/api/save_config", SaveConfigFront)
}
