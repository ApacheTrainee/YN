package router

import (
	"YN/dao"
	"YN/global"
	"YN/log"
	"YN/model"
	"YN/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type fieldFunctionDTO struct {
	Data struct {
		FunctionName string `json:"function_name"`
		Value        bool   `json:"value"`
	} `json:"data"`
}

func PostFieldFunction(c *gin.Context) {
	var req fieldFunctionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WebLogger.Errorf("PostFieldFunction ShouldBindJSON err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "err: Invalid request format",
		})
		return
	}
	log.WebLogger.Infof("%v        %+v", c.Request.URL.Path, req.Data)

	if err := processRCSRequest(req); err != nil {
		log.WebLogger.Errorf("PostFieldFunction processRCSRequest err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("err: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "successful",
	})
}

// 处理RCS请求并写入IO信号
func processRCSRequest(req fieldFunctionDTO) error {
	if isTrue := rasterExclusiveAreaProcess(req); isTrue {
		return nil
	}

	// RCS的bug：有时候会发Value = false的值，这种请求是不要的
	if req.Data.Value == false {
		return nil
	}

	// 就是我们在RCS三方电梯中，自定义配置的函数名 E1_OutReqTo4F
	deviceID := req.Data.FunctionName[:2]   // E1或E2
	signalType := req.Data.FunctionName[3:] // 信号类型(如OutReqTo4F)

	// 安全校验：非agv模式，无法写入
	if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[4] == 0 {
		return fmt.Errorf("device %v is manual mode! Can't operator Elevator", deviceID)
	}

	// 判断电梯是否在线
	if _, exist := global.ElevatorStatus[deviceID]; exist == false {
		return fmt.Errorf("elevator %v disconnect", deviceID)
	}

	// 新增任务
	if signalType == "OutReqTo5F" {
		// 安全校验：不在关门状态，无法新增任务
		if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[1] != 1 {
			return fmt.Errorf("device %v is manual mode! Can't operator Elevator", deviceID)
		}

		// 写入电梯到5楼
		signal := model.ElevatorSignalCoil{
			Write1: 0,
			Write2: 1,
			Write3: 0,
			Write4: 0,
		}
		if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
			return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
		}

		elevatorTask := model.ElevatorTask{
			ElevatorID:  deviceID,
			TaskID:      deviceID + "_" + utils.GenTaskID(),
			StartFloor:  "5F",
			TargetFloor: "4F",
			TaskStatus:  global.ElevatorTaskStatus_To5F_1,
			TaskType:    global.ElevatorTaskType_Down,
			StartTime:   time.Now(),
			UpdateTime:  time.Now(),
		}
		if err := dao.CreateElevatorTask(elevatorTask); err != nil {
			return err
		}

		// 如果原本电梯就是 在5楼、关门、agv模式 状态，直接做逻辑
		if global.ElevatorStatus[deviceID] == 26 {
			// 写入电梯常开门
			signal = model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 1,
				Write4: 0,
			}
			if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
			}

			// 更新任务状态
			elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_Arrive5F_2}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}
		}
	}

	if signalType == "OutReqTo4F" {
		// 安全校验：不在关门状态，无法新增任务
		if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[1] != 1 {
			return fmt.Errorf("device %v is manual mode! Can't operator Elevator", deviceID)
		}

		// 写入电梯到4楼
		signal := model.ElevatorSignalCoil{
			Write1: 1,
			Write2: 0,
			Write3: 0,
			Write4: 0,
		}
		if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
			return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
		}

		// 创建任务
		elevatorTask := model.ElevatorTask{
			ElevatorID:  deviceID,
			TaskID:      deviceID + "_" + utils.GenTaskID(),
			StartFloor:  "4F",
			TargetFloor: "5F",
			TaskStatus:  global.ElevatorTaskStatus_To4F_Up_1,
			TaskType:    global.ElevatorTaskType_Up,
			StartTime:   time.Now(),
			UpdateTime:  time.Now(),
		}
		if err := dao.CreateElevatorTask(elevatorTask); err != nil {
			return err
		}

		if global.ElevatorStatus[deviceID] == 22 {
			// 写入电梯开门
			signal = model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 1,
				Write4: 0,
			}
			if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
			}

			// 更新任务状态
			elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_Arrive4F_Up_2}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}
		}
	}

	// () 查询任务，做状态转变
	elevatorTask, exist := dao.GetElevatorTask(deviceID)
	if exist == false {
		return fmt.Errorf("elevator task not found for elevatorID: %v", deviceID)
	}

	if elevatorTask.TaskType == global.ElevatorTaskType_Down {
		if signalType == "InReqCloseDoor5F" && elevatorTask.TaskStatus == global.ElevatorTaskStatus_OpenDoorFinish_3 {
			// 写入电梯常关门
			signal := model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 0,
				Write4: 1,
			}
			if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
			}

			// 更新任务状态
			elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_CloseDoor_4}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}
		}

		if signalType == "InReqTo4F" {
			if elevatorTask.IsProcessToOtherFloorReq == true {
				return nil
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{IsProcessToOtherFloorReq: true}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}

			i := 0
			for {
				elevatorTask, exist := dao.GetElevatorTask(deviceID)
				if exist == false {
					return fmt.Errorf("elevator task not found for elevatorID: %v", deviceID)
				}

				if elevatorTask.TaskStatus == global.ElevatorTaskStatus_CloseDoorFinish_5 {
					// 写入电梯去4楼
					signal := model.ElevatorSignalCoil{
						Write1: 1,
						Write2: 0,
						Write3: 0,
						Write4: 0,
					}
					if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
						return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
					}

					// 更新任务状态
					elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_To4F_6}
					if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
						return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
					}
				}

				time.Sleep(1 * time.Second)
				i = i + 1

				if i == 120 {
					return nil
				}
			}
		}

		if signalType == "OutReqCloseDoor4F" && elevatorTask.TaskStatus == global.ElevatorTaskStatus_OpenDoorFinish_3 {
			// 写入电梯常关门
			signal := model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 0,
				Write4: 1,
			}
			if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
			}

			// 更新任务状态
			elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_CloseDoor_4}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}
		}
	}

	if elevatorTask.TaskType == global.ElevatorTaskType_Up {
		if signalType == "InReqCloseDoor4F" && elevatorTask.TaskStatus == global.ElevatorTaskStatus_OpenDoorFinish_Up_3 {
			// 写入电梯常关门
			signal := model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 0,
				Write4: 1,
			}
			if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
			}

			// 更新任务状态
			elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_CloseDoor_Up_4}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}
		}

		if signalType == "InReqTo5F" {
			if elevatorTask.IsProcessToOtherFloorReq == true {
				return nil
			}

			// 更新任务状态
			elevatorTaskObj := model.ElevatorTask{IsProcessToOtherFloorReq: true}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTaskObj, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}

			i := 0
			for {
				elevatorTask, exist := dao.GetElevatorTask(deviceID)
				if exist == false {
					return fmt.Errorf("elevator task not found for elevatorID: %v", deviceID)
				}

				if elevatorTask.TaskStatus == global.ElevatorTaskStatus_CloseDoorFinish_Up_5 {
					// 写入电梯去5楼
					signal := model.ElevatorSignalCoil{
						Write1: 0,
						Write2: 1,
						Write3: 0,
						Write4: 0,
					}
					if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
						return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
					}

					// 更新任务状态
					elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_To5F_Up_6}
					if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
						return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
					}
				}

				time.Sleep(1 * time.Second)
				i = i + 1

				if i == 120 {
					return nil
				}
			}
		}

		if signalType == "OutReqCloseDoor5F" && elevatorTask.TaskStatus == global.ElevatorTaskStatus_OpenDoorFinish_Up_3 {
			// 写入电梯常关门
			signal := model.ElevatorSignalCoil{
				Write1: 0,
				Write2: 0,
				Write3: 0,
				Write4: 1,
			}
			if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
			}

			// 更新任务状态
			elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_CloseDoor_Up_4}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}
		}
	}

	return nil
}

func rasterExclusiveAreaProcess(req fieldFunctionDTO) bool {
	if req.Data.FunctionName == "RasterExclusiveArea1" {
		global.RasterExclusiveAreaChan1 <- req.Data.Value
		return true
	}

	if req.Data.FunctionName == "RasterExclusiveArea2" {
		global.RasterExclusiveAreaChan2 <- req.Data.Value
		return true
	}

	return false
}
