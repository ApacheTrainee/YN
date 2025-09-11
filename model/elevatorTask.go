package model

import "time"

type ElevatorTask struct {
	ElevatorID                      string // 电梯ID
	TaskID                          string // 任务ID
	IsProcessToOtherFloorReq        bool
	IsProcessStartFloorCloseDoorReq bool
	StartFloor                      float64 // 开始楼层
	TargetFloor                     float64 // 目标楼层
	TaskStatus                      string  // 任务状态，看全局变量那边
	ReqStatus                       string  // 请求给RCS的状态，用于自动化测试，无业务作用
	TaskType                        string  // 任务类型。up：上楼、down：下楼
	StartTime                       time.Time
	UpdateTime                      time.Time
	EndTime                         time.Time
}
