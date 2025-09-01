package global

import (
	"YN/model"
	"github.com/goburrow/modbus"
	"sync"
)

var (
	TaskLock         sync.RWMutex
	ElevatorTaskList = make(map[string]model.ElevatorTask) //各个电梯的任务列表

	ElevatorStatus = make(map[string]int)           // 读取(叫串口)的电梯状态值
	ClientList     = make(map[string]modbus.Client) // 电梯客户端

	RasterExclusiveAreaChan1 = make(chan bool, 1)
	RasterExclusiveAreaChan2 = make(chan bool, 1)
	RasterExclusiveArea1     = false
	RasterExclusiveArea2     = false
)

var (
	// 电梯任务类型
	ElevatorTaskType_Down = "Down"
	ElevatorTaskType_Up   = "Up"

	// 电梯任务状态流转
	ElevatorTaskStatus_To5F_1            = "To5F"
	ElevatorTaskStatus_Arrive5F_2        = "Arrive5FFinish"
	ElevatorTaskStatus_OpenDoorFinish_3  = "OpenDoorFinish"
	ElevatorTaskStatus_CloseDoor_4       = "CloseDoor"
	ElevatorTaskStatus_CloseDoorFinish_5 = "CloseDoorFinish"
	ElevatorTaskStatus_To4F_6            = "To4F"
	ElevatorTaskStatus_Arrive4F_7        = "Arrive4FFinish"

	ElevatorTaskStatus_To4F_Up_1            = "To4F"
	ElevatorTaskStatus_Arrive4F_Up_2        = "Arrive4FFinish"
	ElevatorTaskStatus_OpenDoorFinish_Up_3  = "OpenDoorFinish"
	ElevatorTaskStatus_CloseDoor_Up_4       = "CloseDoor"
	ElevatorTaskStatus_CloseDoorFinish_Up_5 = "CloseDoorFinish"
	ElevatorTaskStatus_To5F_Up_6            = "To5F"
	ElevatorTaskStatus_Arrive5FFinish_Up_7  = "Arrive5FFinish"

	/*
		web接收的是这8种信号
			E1_OutReqTo5F				AGV在电梯外，申请电梯到5楼
			E1_InReqCloseDoor5F			AGV在电梯内，AGV送货关门请求
			E1_InReqTo4F				AGV在电梯内，申请电梯到4楼
			E1_OutReqCloseDoor4F		AGV在电梯外，AGV送货关门请求
			---------------------
			E1_OutReqTo4F				AGV在电梯外，申请电梯到4楼
			E1_InReqCloseDoor4F			AGV在电梯内，AGV送货关门请求
			E1_InReqTo5F				AGV在电梯内，申请电梯到5楼
			E1_OutReqCloseDoor5F		AGV在电梯外，AGV送货关门请求
	*/
	// 自己发的是这4种信号
	ElevatorRcsConfig_E1_InOpenInPlace5F  = "E1_InOpenInPlace5F"  // AGV在电梯外，请求打开5楼电梯进入
	ElevatorRcsConfig_E1_OutOpenInPlace4F = "E1_OutOpenInPlace4F" // AGV在电梯内，请求打开4楼电梯出来
	ElevatorRcsConfig_E1_InOpenInPlace4F  = "E1_InOpenInPlace4F"  // AGV在电梯外，请求打开4楼电梯进入
	ElevatorRcsConfig_E1_OutOpenInPlace5F = "E1_OutOpenInPlace5F" // AGV在电梯内，请求打开5楼电梯出来
)
