package model

import "time"

type ElevatorTask struct {
	ElevatorID               string // 电梯ID
	TaskID                   string // 任务ID
	IsProcessToOtherFloorReq bool
	StartFloor               string // 开始楼层
	TargetFloor              string // 目标楼层
	TaskStatus               string // 任务状态，看全局变量那边
	TaskType                 string // 任务类型。up：上楼、down：下楼
	StartTime                time.Time
	UpdateTime               time.Time
	EndTime                  time.Time
}
