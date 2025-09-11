package global

import (
	"YN/model"
	"github.com/goburrow/modbus"
	"sync"
)

var (
	TaskLock         sync.RWMutex
	TaskPoolLock     sync.RWMutex
	ElevatorTaskPool = make(map[string][]model.ElevatorTask)
	ElevatorTaskList = make(map[string]model.ElevatorTask) //各个电梯的任务列表

	ElevatorStatus = make(map[string]int)           // 读取(叫串口)的电梯状态值
	ClientList     = make(map[string]modbus.Client) // 电梯客户端

	ElevatorCoilSimulation = []byte{16}

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
	ElevatorTaskStatus_ToStartFloor              = "ToStartFloor"
	ElevatorTaskStatus_StartFloorArriveFinish    = "StartFloorArriveFinish" // 因为生成任务那里要判断是否在当前楼层，所以多2个引用
	ElevatorTaskStatus_StartFloorOpenDoorFinish  = "StartFloorOpenDoorFinish"
	ElevatorTaskStatus_StartFloorCloseDoor       = "StartFloorCloseDoor"
	ElevatorTaskStatus_StartFloorCloseDoorFinish = "StartFloorCloseDoorFinish"

	ElevatorTaskStatus_ToTargetFloor              = "ToTargetFloor"
	ElevatorTaskStatus_TargetFloorArriveFinish    = "TargetFloorArriveFinish"
	ElevatorTaskStatus_TargetFloorOpenDoorFinish  = "TargetFloorOpenDoorFinish"
	ElevatorTaskStatus_TargetFloorCloseDoor       = "TargetFloorCloseDoor"       // 22、26各多两个判断。共多4个
	ElevatorTaskStatus_TargetFloorCloseDoorFinish = "TargetFloorCloseDoorFinish" // 多两个仅更新状态的引用

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
	//ElevatorRcsConfig_E1_OutReqTo5F        = "OutReqTo5F"
	//ElevatorRcsConfig_E1_InReqCloseDoor5F  = "InReqCloseDoor5F"
	//ElevatorRcsConfig_E1_InReqTo4F         = "InReqTo4F"
	//ElevatorRcsConfig_E1_OutReqCloseDoor4F = "OutReqCloseDoor4F"
	//
	//ElevatorRcsConfig_E1_OutReqTo4F        = "OutReqTo4F"
	//ElevatorRcsConfig_E1_InReqCloseDoor4F  = "InReqCloseDoor4F"
	//ElevatorRcsConfig_E1_InReqTo5F         = "InReqTo5F"
	//ElevatorRcsConfig_E1_OutReqCloseDoor5F = "OutReqCloseDoor5F"

	// 自己发的是这4种信号
	ElevatorRcsConfig_E1_InOpenInPlace5F  = "E1_InOpenInPlace5F"  // AGV在电梯外，请求打开5楼电梯进入
	ElevatorRcsConfig_E1_OutOpenInPlace4F = "E1_OutOpenInPlace4F" // AGV在电梯内，请求打开4楼电梯出来
	ElevatorRcsConfig_E1_InOpenInPlace4F  = "E1_InOpenInPlace4F"  // AGV在电梯外，请求打开4楼电梯进入
	ElevatorRcsConfig_E1_OutOpenInPlace5F = "E1_OutOpenInPlace5F" // AGV在电梯内，请求打开5楼电梯出来
)
