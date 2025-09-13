package service

import (
	"YN/config"
	"YN/dao"
	"YN/global"
	"YN/log"
	"YN/model"
	"YN/utils"
	"fmt"
	"math"
	"sync"
	"time"
)

func ElevatorTaskPoolProcess() {
	var wg sync.WaitGroup
	for _, device := range config.Config.DeviceList {
		wg.Add(1)
		go func(device config.Device) {
			defer wg.Done()

			elevatorTaskPoolProcessImp(device.Id)
		}(device)
	}

	wg.Wait()
}

func elevatorTaskPoolProcessImp(deviceID string) {
	for {
		time.Sleep(1 * time.Second)

		if len(global.ElevatorTaskPool[deviceID]) == 0 {
			continue
		}
		// 如果正在执行任务，也continue
		if global.ElevatorTaskList[deviceID].TaskStatus != "" {
			continue
		}

		var currentFloor float64 = 5
		if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[2] == 1 {
			currentFloor = 4 // 在四楼
		}

		// ()判断电梯当前位置，绝对值最小那个任务，如果有多个，取创建时间最早那个
		firstTask := global.ElevatorTaskPool[deviceID][0]
		minFloor := math.Abs(firstTask.StartFloor - currentFloor)
		elevatorTask := firstTask

		// 挑选合适的电梯任务：电梯移动层数最少的任务。都相同就取创建时间早的任务
		global.TaskPoolLock.Lock()
		for _, elevatorTaskObj := range global.ElevatorTaskPool[deviceID] {
			currentDiff := math.Abs(elevatorTaskObj.StartFloor - currentFloor)

			if currentDiff < minFloor {
				minFloor = currentDiff
				elevatorTask = elevatorTaskObj
			} else if minFloor == elevatorTask.StartFloor {
				if elevatorTaskObj.StartTime.Before(elevatorTask.StartTime) {
					elevatorTask = elevatorTaskObj
				}
			}
		}
		global.TaskPoolLock.Unlock()

		// 删除任务池中的任务
		global.TaskPoolLock.Lock()
		index := -1
		for i, elevatorTaskObj := range global.ElevatorTaskPool[deviceID] {
			if elevatorTaskObj.TaskID == elevatorTask.TaskID {
				index = i
				break
			}
		}
		if index != -1 {
			global.ElevatorTaskPool[deviceID] = append(global.ElevatorTaskPool[deviceID][:index], global.ElevatorTaskPool[deviceID][index+1:]...)
		}
		global.TaskPoolLock.Unlock()

		// 合适的电梯任务，替换全局电梯任务
		global.ElevatorTaskList[deviceID] = elevatorTask

		// 触发任务处理逻辑
		if err := taskSchedule(currentFloor, global.ElevatorTaskList[deviceID], deviceID); err != nil {
			if err = taskSchedule(currentFloor, global.ElevatorTaskList[deviceID], deviceID); err != nil { // 失败了再试一次
				log.Logger.Errorf("err: %v", err) // 最后不行，也只能打日志，停下来了【只能取消任务，重新再跑了】
			}
		}
	}
}

func taskSchedule(currentFloor float64, elevatorTask model.ElevatorTask, deviceID string) error {
	currentDoorStatus := false // 关门【默认】
	if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[0] == 1 {
		currentDoorStatus = true
	}

	// 电梯在其他楼层
	if elevatorTask.StartFloor != currentFloor {
		if currentDoorStatus == false { // 如果是关着门的
			if elevatorTask.StartFloor == 5 {
				// 写入电梯到5楼
				if err := utils.WriteElevatorCoils(deviceID, model.ElevatorSignalCoil{ReqTo5F2: 1}); err != nil {
					return fmt.Errorf("elevator = %v writeData error: %v", deviceID, err)
				}
			}

			if elevatorTask.StartFloor == 4 {
				// 写入电梯到4楼
				if err := utils.WriteElevatorCoils(deviceID, model.ElevatorSignalCoil{ReqTo4F1: 1}); err != nil {
					return fmt.Errorf("elevator = %v writeData error: %v", deviceID, err)
				}
			}
		}

		if currentDoorStatus == true { // 如果门是开着的
			// 写入电梯关门
			if err := utils.WriteElevatorCoils(deviceID, model.ElevatorSignalCoil{CloseDoor4: 1}); err != nil {
				return fmt.Errorf("elevator = %v writeData error: %v", deviceID, err)
			}
		}
	}

	// 电梯在当前楼层
	if elevatorTask.StartFloor == currentFloor {
		if currentDoorStatus == false {
			// 写入电梯常开门
			if err := utils.WriteElevatorCoils(deviceID, model.ElevatorSignalCoil{OpenDoor3: 1}); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
			}

			// 更新任务状态
			elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_StartFloorArriveFinish}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}
		}

		if currentDoorStatus == true { // 如果开门到位
			var elevatorTaskObj model.ElevatorTask

			if currentFloor == 5 {
				if config.Config.RunMode == "pro" {
					// 发送RCS
					if err := utils.SendRCS(global.ElevatorRcsConfig_E1_InOpenInPlace5F, true); err != nil {
						return fmt.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
					}
				} else {
					log.Logger.Infof("send to RCS: request body is %v", global.ElevatorRcsConfig_E1_InOpenInPlace5F)
				}

				elevatorTaskObj = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_StartFloorOpenDoorFinish, ReqStatus: global.ElevatorRcsConfig_E1_InOpenInPlace5F}
			}
			if currentFloor == 4 {
				if config.Config.RunMode == "pro" {
					// 发送RCS
					if err := utils.SendRCS(global.ElevatorRcsConfig_E1_InOpenInPlace4F, true); err != nil {
						return fmt.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
					}
				} else {
					log.Logger.Infof("send to RCS: request body is %v", global.ElevatorRcsConfig_E1_InOpenInPlace4F)
				}

				elevatorTaskObj = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_StartFloorOpenDoorFinish, ReqStatus: global.ElevatorRcsConfig_E1_InOpenInPlace4F}
			}

			// 更新任务状态
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}
	}

	return nil
}
