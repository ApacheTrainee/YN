package service

import (
	"YN/config"
	"YN/dao"
	"YN/global"
	"YN/log"
	"YN/model"
	"YN/utils"
	"fmt"
	"time"
)

// 读取到电梯状态变化后，处理电梯信号变化
func processElevatorSignal(device config.Device) {
	// 校验：判断是否为进行中的任务
	elevatorTask, exist := dao.GetElevatorTask(device.Id)
	if exist == false {
		return
	}

	// 校验：是否为AGV模式
	if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[device.Id])}); bits[4] != 1 {
		log.Logger.Errorf("not AgvAutoMode: %v", device.Id)
		return
	}

	// 电梯当前在哪楼
	// 跨楼层时，哪层都不在，都是0
	var currentFloor float64 = 0
	if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[device.Id])}); bits[3] == 1 {
		currentFloor = 5 // 在四楼
	}
	if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[device.Id])}); bits[2] == 1 {
		currentFloor = 4 // 在四楼
	}

	// 电梯当前开门还是关门。用bool类型不行，因为开关中都是0，会造成默认是关门的
	var currentDoorStatus string
	if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[device.Id])}); bits[0] == 1 {
		currentDoorStatus = "open"
	}
	if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[device.Id])}); bits[1] == 1 {
		currentDoorStatus = "close"
	}

	// 根据任务状态做判断
	if err := elevatorTaskStatusForToStartFloorProcess(elevatorTask, currentDoorStatus, currentFloor, device); err != nil {
		log.Logger.Errorf("%v", err)
		return
	}
	if err := elevatorTaskStatusForStartFloorArriveFinishProcess(elevatorTask, currentDoorStatus, currentFloor, device); err != nil {
		log.Logger.Errorf("%v", err)
		return
	}
	if err := elevatorTaskStatusForStartFloorCloseDoorProcess(elevatorTask, currentDoorStatus, currentFloor, device); err != nil {
		log.Logger.Errorf("%v", err)
		return
	}

	if err := elevatorTaskStatusForToTargetFloorProcess(elevatorTask, currentDoorStatus, currentFloor, device); err != nil {
		log.Logger.Errorf("%v", err)
		return
	}
	if err := elevatorTaskStatusForTargetFloorArriveFinishProcess(elevatorTask, currentDoorStatus, currentFloor, device); err != nil {
		log.Logger.Errorf("%v", err)
		return
	}
	if err := elevatorTaskStatusForTargetFloorCloseDoorProcess(elevatorTask, currentDoorStatus, currentFloor, device); err != nil {
		log.Logger.Errorf("%v", err)
		return
	}
}

func elevatorTaskStatusForToStartFloorProcess(elevatorTask model.ElevatorTask, currentDoorStatus string, currentFloor float64, device config.Device) error {
	if elevatorTask.TaskStatus == global.ElevatorTaskStatus_ToStartFloor {
		if elevatorTask.StartFloor != currentFloor {
			if currentDoorStatus == "close" { // 如果是关着门的
				if elevatorTask.StartFloor == 5 {
					// 写入电梯到5楼
					if err := utils.WriteElevatorCoils(device.Id, model.ElevatorSignalCoil{ReqTo5F2: 1}); err != nil {
						return fmt.Errorf("elevator = %v writeData error: %v", device.Id, err)
					}
				}

				if elevatorTask.StartFloor == 4 {
					// 写入电梯到4楼
					if err := utils.WriteElevatorCoils(device.Id, model.ElevatorSignalCoil{ReqTo4F1: 1}); err != nil {
						return fmt.Errorf("elevator = %v writeData error: %v", device.Id, err)
					}
				}
			}

			// 门开着的处理，已经另一边处理过了
		}

		// 在当前楼层的处理
		if elevatorTask.StartFloor == currentFloor {
			// 门关着的，就打开
			if currentDoorStatus == "close" {
				// 写入电梯开门
				if err := utils.WriteElevatorCoils(elevatorTask.ElevatorID, model.ElevatorSignalCoil{OpenDoor3: 1}); err != nil {
					return fmt.Errorf("deviceid = %v reset writeData error: %v", elevatorTask.ElevatorID, err)
				}

				// 更新任务状态
				elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_StartFloorArriveFinish}
				if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
					return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
				}
			}

			// 门开着的情况，初始任务调度那边处理了
		}
	}

	return nil
}

func elevatorTaskStatusForStartFloorArriveFinishProcess(elevatorTask model.ElevatorTask, currentDoorStatus string, currentFloor float64, device config.Device) error {
	//
	if elevatorTask.TaskStatus == global.ElevatorTaskStatus_StartFloorArriveFinish {
		var elevatorTaskObj model.ElevatorTask

		if currentDoorStatus == "open" { // 如果开门到位
			if currentFloor == 5 {
				if config.Config.RunMode == "pro" {
					// 发送RCS
					if err := utils.SendRCS(global.ElevatorRcsConfig_E1_InOpenInPlace5F, true); err != nil {
						return fmt.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
					}
				} else {
					log.Logger.Infof("----6666-----send to RCS: request body is %v", global.ElevatorRcsConfig_E1_InOpenInPlace5F)
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
					log.Logger.Infof("----6666-----send to RCS: request body is %v", global.ElevatorRcsConfig_E1_InOpenInPlace4F)
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

func elevatorTaskStatusForStartFloorCloseDoorProcess(elevatorTask model.ElevatorTask, currentDoorStatus string, currentFloor float64, device config.Device) error {
	if elevatorTask.TaskStatus == global.ElevatorTaskStatus_StartFloorCloseDoor {
		if currentDoorStatus == "close" {
			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_StartFloorCloseDoorFinish}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}
	}

	return nil
}

func elevatorTaskStatusForToTargetFloorProcess(elevatorTask model.ElevatorTask, currentDoorStatus string, currentFloor float64, device config.Device) error {
	if elevatorTask.TaskStatus == global.ElevatorTaskStatus_ToTargetFloor {
		if currentDoorStatus == "close" && elevatorTask.TargetFloor == currentFloor {
			// 写入电梯常开门
			if err := utils.WriteElevatorCoils(elevatorTask.ElevatorID, model.ElevatorSignalCoil{OpenDoor3: 1}); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", elevatorTask.ElevatorID, err)
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_TargetFloorArriveFinish}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}
	}

	return nil
}

func elevatorTaskStatusForTargetFloorArriveFinishProcess(elevatorTask model.ElevatorTask, currentDoorStatus string, currentFloor float64, device config.Device) error {
	if elevatorTask.TaskStatus == global.ElevatorTaskStatus_TargetFloorArriveFinish {
		if currentDoorStatus == "open" && elevatorTask.TargetFloor == currentFloor {
			var elevatorTaskObj model.ElevatorTask
			if currentFloor == 5 {
				if config.Config.RunMode == "pro" {
					// 发送RCS
					if err := utils.SendRCS(global.ElevatorRcsConfig_E1_OutOpenInPlace5F, true); err != nil {
						return fmt.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
					}
				} else {
					log.Logger.Infof("----6666-----send to RCS: request body is %v", global.ElevatorRcsConfig_E1_OutOpenInPlace5F)
				}

				elevatorTaskObj = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_TargetFloorOpenDoorFinish, ReqStatus: global.ElevatorRcsConfig_E1_OutOpenInPlace5F}
			}
			if currentFloor == 4 {
				if config.Config.RunMode == "pro" {
					// 发送RCS
					if err := utils.SendRCS(global.ElevatorRcsConfig_E1_OutOpenInPlace4F, true); err != nil {
						return fmt.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
					}
				} else {
					log.Logger.Infof("----6666-----send to RCS: request body is %v", global.ElevatorRcsConfig_E1_OutOpenInPlace4F)
				}

				elevatorTaskObj = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_TargetFloorOpenDoorFinish, ReqStatus: global.ElevatorRcsConfig_E1_OutOpenInPlace4F}
			}

			// 更新任务状态
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}
	}

	return nil
}

func elevatorTaskStatusForTargetFloorCloseDoorProcess(elevatorTask model.ElevatorTask, currentDoorStatus string, currentFloor float64, device config.Device) error {
	if elevatorTask.TaskStatus == global.ElevatorTaskStatus_TargetFloorCloseDoor {
		if currentDoorStatus == "close" && elevatorTask.TargetFloor == currentFloor {
			// 清空电梯写入
			if err := utils.WriteElevatorCoils(elevatorTask.ElevatorID, model.ElevatorSignalCoil{}); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", elevatorTask.ElevatorID, err)
			}

			// 更新任务状态 - 任务完成
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_TargetFloorCloseDoorFinish, EndTime: time.Now()}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, true); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}

			// 删除任务
			dao.DeleteElevatorTask(elevatorTask.ElevatorID)
		}
	}

	return nil
}
