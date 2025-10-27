package service

import (
	"YN/config"
	"YN/global"
	"YN/log"
	"YN/utils"
	"time"
)

// 读取到电梯状态变化后，处理电梯信号变化
func processElevatorSignal(device config.Device) {
	// 校验：读取电梯的值是否正确
	bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[device.Id])})
	if (bits[0] == 0 && bits[1] == 0 && bits[2] == 0 && bits[3] == 0) || (bits[0] == 1 && bits[1] == 1 && bits[2] == 1 && bits[3] == 1) || bits[4] == 0 || bits[5] == 1 {
		log.Logger.Errorf("elevator err. bits: %v", bits)
		return
	}

	// 电梯当前在哪楼
	// 跨楼层时，哪层都不在，都是0
	var currentFloor float64 = 0
	if bits[2] == 1 {
		currentFloor = 4 // 在四楼
	}
	if bits[3] == 1 {
		currentFloor = 5 // 在四楼
	}

	// 电梯当前开门还是关门。用bool类型不行，因为开关中都是0，会造成默认是关门的
	currentDoorStatus := ""
	if bits[0] == 1 {
		currentDoorStatus = global.OpenDoor
	}
	if bits[1] == 1 {
		currentDoorStatus = global.CloseDoor
	}

	if currentFloor == 5 && currentDoorStatus == global.CloseDoor {
		if global.ElevatorTask.Status == global.ElevatorTaskStatus_ToStartFloor && global.ElevatorTask.StartFloor == currentFloor {
			// 5楼电梯开门
			if err := utils.WriteElevatorCoils(device.Id, 0, global.OpenDoor); err != nil {
				log.Logger.Errorf("elevator = %v writeData error: %v", device.Id, err)
				return
			}

			global.ElevatorTask.Status = global.ElevatorTaskStatus_StartFloorArriveFinish
			global.ElevatorTask.StartTime = time.Now()
		}

		if global.ElevatorTask.Status == global.ElevatorTaskStatus_ToTargetFloor && global.ElevatorTask.TargetFloor == currentFloor {
			// 5楼电梯开门
			if err := utils.WriteElevatorCoils(device.Id, 0, global.OpenDoor); err != nil {
				log.Logger.Errorf("elevator = %v writeData error: %v", device.Id, err)
				return
			}

			global.ElevatorTask.Status = global.ElevatorTaskStatus_TargetFloorArriveFinish
			global.ElevatorTask.StartTime = time.Now()
		}
	}
	if currentFloor == 5 && currentDoorStatus == global.OpenDoor {
		if global.ElevatorTask.Status == global.ElevatorTaskStatus_StartFloorArriveFinish {
			global.ElevatorTask.Status = ""

			// agv进5楼电梯
			if config.Config.RunMode == "pro" {
				// 发送RCS
				if err := utils.SendRCS(global.ElevatorRcsConfig_E1_InOpenInPlace5F, true); err != nil {
					log.Logger.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
				}
			} else {
				log.Logger.Infof("----4-----send to RCS: request body is %v", global.ElevatorRcsConfig_E1_InOpenInPlace5F)
			}
		}

		if global.ElevatorTask.Status == global.ElevatorTaskStatus_TargetFloorArriveFinish {
			global.ElevatorTask.Status = ""

			// agv出5楼电梯
			if config.Config.RunMode == "pro" {
				// 发送RCS
				if err := utils.SendRCS(global.ElevatorRcsConfig_E1_OutOpenInPlace5F, true); err != nil {
					log.Logger.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
				}
			} else {
				log.Logger.Infof("----3-----send to RCS: request body is %v", global.ElevatorRcsConfig_E1_OutOpenInPlace5F)
			}
		}
	}

	if currentFloor == 4 && currentDoorStatus == global.CloseDoor {
		if global.ElevatorTask.Status == global.ElevatorTaskStatus_ToStartFloor && global.ElevatorTask.StartFloor == currentFloor {
			// 4楼电梯开门
			if err := utils.WriteElevatorCoils(device.Id, 0, global.OpenDoor); err != nil {
				log.Logger.Errorf("elevator = %v writeData error: %v", device.Id, err)
				return
			}

			global.ElevatorTask.Status = global.ElevatorTaskStatus_StartFloorArriveFinish
			global.ElevatorTask.StartTime = time.Now()
		}

		if global.ElevatorTask.Status == global.ElevatorTaskStatus_ToTargetFloor && global.ElevatorTask.TargetFloor == currentFloor {
			// 4楼电梯开门
			if err := utils.WriteElevatorCoils(device.Id, 0, global.OpenDoor); err != nil {
				log.Logger.Errorf("elevator = %v writeData error: %v", device.Id, err)
				return
			}

			global.ElevatorTask.Status = global.ElevatorTaskStatus_TargetFloorArriveFinish
			global.ElevatorTask.StartTime = time.Now()
		}
	}
	if currentFloor == 4 && currentDoorStatus == global.OpenDoor {
		if global.ElevatorTask.Status == global.ElevatorTaskStatus_StartFloorArriveFinish {
			global.ElevatorTask.Status = ""

			// agv进4楼电梯
			if config.Config.RunMode == "pro" {
				// 发送RCS
				if err := utils.SendRCS(global.ElevatorRcsConfig_E1_InOpenInPlace4F, true); err != nil {
					log.Logger.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
				}
			} else {
				log.Logger.Infof("----1-----send to RCS: request body is %v", global.ElevatorRcsConfig_E1_InOpenInPlace4F)
			}
		}

		if global.ElevatorTask.Status == global.ElevatorTaskStatus_TargetFloorArriveFinish {
			global.ElevatorTask.Status = ""

			// agv出4楼电梯
			if config.Config.RunMode == "pro" {
				// 发送RCS
				if err := utils.SendRCS(global.ElevatorRcsConfig_E1_OutOpenInPlace4F, true); err != nil {
					log.Logger.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
				}
			} else {
				log.Logger.Infof("----2-----send to RCS: request body is %v", global.ElevatorRcsConfig_E1_OutOpenInPlace4F)
			}
		}
	}
}
