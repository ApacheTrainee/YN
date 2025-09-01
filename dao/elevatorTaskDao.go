package dao

import (
	"YN/global"
	"YN/log"
	"YN/model"
	"fmt"
	"time"
)

func CreateElevatorTask(elevatorTask model.ElevatorTask) error {
	global.TaskLock.Lock()
	defer global.TaskLock.Unlock()

	// 检查是否已存在该电梯的任务
	if _, exists := global.ElevatorTaskList[elevatorTask.ElevatorID]; exists {
		return fmt.Errorf("---- GenTask Fail！------, deviceId:%v already existsTask, system Can't genTask", elevatorTask.ElevatorID) // 已存在，不覆盖
	}

	global.ElevatorTaskList[elevatorTask.ElevatorID] = elevatorTask

	log.Logger.Infof("create ElevatorTaskList: %+v", global.ElevatorTaskList[elevatorTask.ElevatorID])
	return nil
}

func GetElevatorTask(elevatorID string) (model.ElevatorTask, bool) {
	global.TaskLock.RLock()
	defer global.TaskLock.RUnlock()

	task, exist := global.ElevatorTaskList[elevatorID]
	return task, exist
}

func DeleteElevatorTask(elevatorID string) {
	global.TaskLock.Lock()
	defer global.TaskLock.Unlock()

	delete(global.ElevatorTaskList, elevatorID)
	log.Logger.Infof("delete ElevatorTaskList: %+v\n", elevatorID)
}

func UpdateElevatorTask(elevatorID string, elevatorTask model.ElevatorTask, isFinish bool) error {
	global.TaskLock.Lock()
	defer global.TaskLock.Unlock()

	task, exists := global.ElevatorTaskList[elevatorID]
	if !exists {
		return fmt.Errorf("elevator task not found for elevatorID: %v", elevatorID)
	}

	// 链式判断更新
	if elevatorTask.ElevatorID != "" {
		task.ElevatorID = elevatorTask.ElevatorID
	}
	if elevatorTask.TaskID != "" {
		task.TaskID = elevatorTask.TaskID
	}
	if elevatorTask.IsProcessToOtherFloorReq != false {
		task.IsProcessToOtherFloorReq = elevatorTask.IsProcessToOtherFloorReq
	}
	if elevatorTask.StartFloor != "" {
		task.StartFloor = elevatorTask.StartFloor
	}
	if elevatorTask.TargetFloor != "" {
		task.TargetFloor = elevatorTask.TargetFloor
	}
	if elevatorTask.TaskStatus != "" {
		task.TaskStatus = elevatorTask.TaskStatus
	}
	if elevatorTask.TaskType != "" {
		task.TaskType = elevatorTask.TaskType
	}
	if isFinish == true {
		task.EndTime = elevatorTask.EndTime
	}
	task.UpdateTime = time.Now()

	global.ElevatorTaskList[elevatorID] = task

	log.Logger.Infof("update ElevatorTaskList-now: %+v", global.ElevatorTaskList[elevatorID])
	return nil
}
