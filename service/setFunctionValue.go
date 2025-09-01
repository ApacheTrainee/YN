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
func processElevatorSignal(device config.Device, signalValue int) {
	// 校验：判断是否为进行中的任务
	elevatorTask, exist := dao.GetElevatorTask(device.Id)
	if exist == false {
		return
	}

	// 校验：是否为AGV模式
	if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[device.Id])}); bits[4] != 1 {
		log.Logger.Errorf("not AutoMode: %v", device.Id)
		return
	}

	/*
		读取电梯状态，用到的状态：
		101010 [21]	4楼电梯开门到位
		100110 [25]	5楼电梯开门到位
		011010 [22] 原4楼电梯关门、电梯达到4楼、4楼电梯关门到位
		010110 [26]	原5楼电梯关门、电梯到达5楼、5楼电梯关门到位
	*/
	if err := processSignalFor21(elevatorTask, signalValue); err != nil {
		log.Logger.Errorf("%v", err)
		return
	}
	if err := processSignalFor25(elevatorTask, signalValue); err != nil {
		log.Logger.Errorf("%v", err)
		return
	}
	if err := processSignalFor22(elevatorTask, signalValue); err != nil {
		log.Logger.Errorf("%v", err)
		return
	}
	if err := processSignalFor26(elevatorTask, signalValue); err != nil {
		log.Logger.Errorf("%v", err)
		return
	}
}

func processSignalFor21(elevatorTask model.ElevatorTask, signalValue int) error {
	if signalValue != 21 {
		return nil
	}

	if elevatorTask.TaskType == global.ElevatorTaskType_Down {
		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_Arrive4F_7 {
			if config.Config.RunMode == "pro" {
				// 发送RCS
				if err := utils.SendRCS(global.ElevatorRcsConfig_E1_OutOpenInPlace4F, true); err != nil {
					return fmt.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
				}
			} else {
				log.Logger.Infof("send to RCS: request body is %v", global.ElevatorRcsConfig_E1_OutOpenInPlace4F)
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_OpenDoorFinish_3}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}
	}

	if elevatorTask.TaskType == global.ElevatorTaskType_Up {
		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_Arrive4F_Up_2 {
			if config.Config.RunMode == "pro" {
				// 发送RCS
				if err := utils.SendRCS(global.ElevatorRcsConfig_E1_InOpenInPlace4F, true); err != nil {
					return fmt.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
				}
			} else {
				log.Logger.Infof("send to RCS: request body is %v", global.ElevatorRcsConfig_E1_InOpenInPlace4F)
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_OpenDoorFinish_Up_3}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}
	}

	return nil
}

func processSignalFor25(elevatorTask model.ElevatorTask, signalValue int) error {
	if signalValue != 25 {
		return nil
	}

	if elevatorTask.TaskType == global.ElevatorTaskType_Down {
		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_Arrive5F_2 {
			if config.Config.RunMode == "pro" {
				// 发送RCS
				if err := utils.SendRCS(global.ElevatorRcsConfig_E1_InOpenInPlace5F, true); err != nil {
					return fmt.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
				}
			} else {
				log.Logger.Infof("send to RCS: request body is %v", global.ElevatorRcsConfig_E1_InOpenInPlace5F)
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_OpenDoorFinish_3}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}
	}

	if elevatorTask.TaskType == global.ElevatorTaskType_Up {
		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_Arrive5FFinish_Up_7 {
			if config.Config.RunMode == "pro" {
				// 发送RCS
				if err := utils.SendRCS(global.ElevatorRcsConfig_E1_OutOpenInPlace5F, true); err != nil {
					return fmt.Errorf("SendRCS updateRCSFieldFunction err: %v", err)
				}
			} else {
				log.Logger.Infof("send to RCS: request body is %v", global.ElevatorRcsConfig_E1_OutOpenInPlace5F)
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_OpenDoorFinish_Up_3}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}
	}

	return nil
}

func processSignalFor22(elevatorTask model.ElevatorTask, signalValue int) error {
	if signalValue != 22 {
		return nil
	}

	// 注意这里有多个状态需要分开处理
	if elevatorTask.TaskType == global.ElevatorTaskType_Down {
		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_To4F_6 {
			// 写入电梯常开门
			signal := model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 1,
				Write4: 0,
			}
			if err := utils.WriteElevatorCoils(elevatorTask.ElevatorID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", elevatorTask.ElevatorID, err)
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_Arrive4F_7}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}

		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_CloseDoor_4 {
			// 清空电梯写入
			signal := model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 0,
				Write4: 0,
			}
			if err := utils.WriteElevatorCoils(elevatorTask.ElevatorID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", elevatorTask.ElevatorID, err)
			}

			// 更新任务状态 - 任务完成
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_CloseDoorFinish_5, EndTime: time.Now()}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, true); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}

			// 删除任务
			dao.DeleteElevatorTask(elevatorTask.ElevatorID)
		}
	}

	if elevatorTask.TaskType == global.ElevatorTaskType_Up {
		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_To4F_Up_1 {
			// 写入电梯开门
			signal := model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 1,
				Write4: 0,
			}
			if err := utils.WriteElevatorCoils(elevatorTask.ElevatorID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", elevatorTask.ElevatorID, err)
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_Arrive4F_Up_2}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}

		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_CloseDoor_Up_4 {
			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_CloseDoorFinish_Up_5}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}
	}

	return nil
}

func processSignalFor26(elevatorTask model.ElevatorTask, signalValue int) error {
	if signalValue != 26 {
		return nil
	}

	// 注意这里有多个状态需要分开处理
	if elevatorTask.TaskType == global.ElevatorTaskType_Down {
		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_To5F_1 {
			// 写入电梯常开门
			signal := model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 1,
				Write4: 0,
			}
			if err := utils.WriteElevatorCoils(elevatorTask.ElevatorID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", elevatorTask.ElevatorID, err)
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_Arrive5F_2}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}

		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_CloseDoor_4 {
			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_CloseDoorFinish_5}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}

	}

	if elevatorTask.TaskType == global.ElevatorTaskType_Up {
		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_To5F_Up_6 {
			// 写入电梯常开门
			signal := model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 1,
				Write4: 0,
			}
			if err := utils.WriteElevatorCoils(elevatorTask.ElevatorID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", elevatorTask.ElevatorID, err)
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_Arrive5FFinish_Up_7}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}
		}

		if elevatorTask.TaskStatus == global.ElevatorTaskStatus_CloseDoor_Up_4 {
			// 清空电梯写入
			signal := model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 0,
				Write4: 0,
			}
			if err := utils.WriteElevatorCoils(elevatorTask.ElevatorID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", elevatorTask.ElevatorID, err)
			}

			// 更新任务状态 - 任务完成
			elevatorTaskObj := model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_CloseDoorFinish_Up_5, EndTime: time.Now()}
			if err := dao.UpdateElevatorTask(elevatorTask.ElevatorID, elevatorTaskObj, true); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", elevatorTask.ElevatorID, elevatorTask, global.ElevatorTaskList[elevatorTask.ElevatorID])
			}

			// 删除任务
			dao.DeleteElevatorTask(elevatorTask.ElevatorID)
		}
	}

	return nil
}
