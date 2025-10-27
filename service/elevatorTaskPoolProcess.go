package service

import (
	"YN/config"
	"YN/global"
	"YN/log"
	"YN/utils"
	"fmt"
	"time"
)

func ElevatorTaskPoolProcess() {
	for _, device := range config.Config.DeviceList {
		go func(device config.Device) {
			elevatorTaskPoolProcessImp(device.Id)
		}(device)
	}
}

func elevatorTaskPoolProcessImp(deviceID string) {
	// 遍历管道会一直阻塞，直到有数据进来
	for startFloor := range global.StartFloorProcessChan {
		// 电梯当前在哪楼
		var currentFloor float64 = 0
		if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[3] == 1 {
			currentFloor = 5 // 在四楼
		}
		if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[2] == 1 {
			currentFloor = 4 // 在四楼
		}
		if currentFloor == 0 {
			log.Logger.Errorf("elevatorTaskPoolProcessImp, currentFloor = %v, deviceID: %v", currentFloor, deviceID)
			global.StartFloorProcessChanResult <- "ng"
			continue
		}

		// 电梯当前开门还是关门。用bool类型不行，因为开关中都是0，会造成默认是关门的
		currentDoorStatus := ""
		if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[0] == 1 {
			currentDoorStatus = global.OpenDoor
		}
		if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[1] == 1 {
			currentDoorStatus = global.CloseDoor
		}
		if currentDoorStatus == "" { // 有可能门正好处于关门中
			doorStatus, err := readDoorStatus(deviceID, currentDoorStatus) // 给10秒时间重复读取状态
			if err != nil {
				log.Logger.Errorf("elevatorTaskPoolProcessImp, readDoorStatus error: %v, deviceID: %v", err, deviceID)
				global.StartFloorProcessChanResult <- "ng"
				continue
			}

			currentDoorStatus = doorStatus
		}

		// 电梯在其他楼层
		if startFloor != currentFloor {
			// 写入到起始楼层 和 电梯关门
			if err := utils.WriteElevatorCoils(deviceID, startFloor, global.CloseDoor); err != nil {
				log.Logger.Errorf("elevator = %v writeData error: %v", deviceID, err)
				global.StartFloorProcessChanResult <- "ng"
				continue
			}
		}

		// 电梯在当前楼层
		if startFloor == currentFloor {
			if currentDoorStatus == global.CloseDoor {
				// 写入电梯常开门
				if err := utils.WriteElevatorCoils(deviceID, 0, global.OpenDoor); err != nil {
					log.Logger.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
					global.StartFloorProcessChanResult <- "ng"
					continue
				}

				global.ElevatorTask.Status = global.ElevatorTaskStatus_StartFloorArriveFinish
			}

			if currentDoorStatus == global.OpenDoor { // 如果开门到位
				global.ElevatorTask.Status = ""

				if currentFloor == 5 {
					if config.Config.RunMode == "pro" {
						// 发送RCS
						if err := utils.SendRCS(global.ElevatorRcsConfig_E1_InOpenInPlace5F, true); err != nil {
							log.Logger.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
							global.StartFloorProcessChanResult <- "ng"
							continue
						}
					} else {
						log.Logger.Infof("send to RCS: request body is %v", global.ElevatorRcsConfig_E1_InOpenInPlace5F)
					}
				}

				if currentFloor == 4 {
					global.ElevatorTask.Status = ""

					if config.Config.RunMode == "pro" {
						// 发送RCS
						if err := utils.SendRCS(global.ElevatorRcsConfig_E1_InOpenInPlace4F, true); err != nil {
							log.Logger.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
							global.StartFloorProcessChanResult <- "ng"
							continue
						}
					} else {
						log.Logger.Infof("send to RCS: request body is %v", global.ElevatorRcsConfig_E1_InOpenInPlace4F)
					}
				}
			}
		}

		global.StartFloorProcessChanResult <- "ok"
	}
}

func readDoorStatus(deviceID string, currentDoorStatus string) (string, error) {
	i := 0
	for {
		if i >= 10 {
			return "", fmt.Errorf("10秒内，电梯未读取到门状态")
		}
		time.Sleep(1 * time.Second)
		i = i + 1

		if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[0] == 1 {
			currentDoorStatus = global.OpenDoor
		}
		if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[1] == 1 {
			currentDoorStatus = global.CloseDoor
		}
		log.Logger.Infof("读取电梯状态为 %v", currentDoorStatus)

		if currentDoorStatus != "" {
			return currentDoorStatus, nil
		}
	}
}
