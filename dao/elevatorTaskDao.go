package dao

import (
	"YN/global"
	"YN/log"
	"YN/model"
	"fmt"
	"time"
)

func AddElevatorTaskPool(elevatorTask model.ElevatorTask) error {
	global.TaskPoolLock.Lock()
	defer global.TaskPoolLock.Unlock()

	// 检查是否已存在该电梯的任务
	for _, elevatorTaskObj := range global.ElevatorTaskPool[elevatorTask.ElevatorID] {
		if elevatorTaskObj.TaskID == elevatorTask.TaskID {
			return fmt.Errorf("---- GenTask Fail！------, deviceId:%v already existsTask, system Can't genTask", elevatorTask.ElevatorID)
		}
	}

	global.ElevatorTaskPool[elevatorTask.ElevatorID] = append(global.ElevatorTaskPool[elevatorTask.ElevatorID], elevatorTask)

	log.Logger.Infof("create ElevatorTaskPool: %+v", global.ElevatorTaskPool[elevatorTask.ElevatorID])
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
	log.Logger.Infof("delete ElevatorTaskList: %+v", elevatorID)
	log.Logger.Infof("---------------------------------------------------------------------------------\n\n")
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
	if elevatorTask.IsProcessStartFloorCloseDoorReq != false {
		task.IsProcessStartFloorCloseDoorReq = elevatorTask.IsProcessStartFloorCloseDoorReq
	}
	if elevatorTask.StartFloor != 0 {
		task.StartFloor = elevatorTask.StartFloor
	}
	if elevatorTask.TargetFloor != 0 {
		task.TargetFloor = elevatorTask.TargetFloor
	}
	if elevatorTask.TaskStatus != "" {
		task.TaskStatus = elevatorTask.TaskStatus
	}
	if elevatorTask.TaskType != "" {
		task.TaskType = elevatorTask.TaskType
	}
	if elevatorTask.ReqStatus != "" {
		task.ReqStatus = elevatorTask.ReqStatus
	}
	if isFinish == true {
		task.EndTime = elevatorTask.EndTime
	}
	task.UpdateTime = time.Now()

	global.ElevatorTaskList[elevatorID] = task

	log.Logger.Infof("update ElevatorTaskList-now: %+v", global.ElevatorTaskList[elevatorID])
	return nil
}
