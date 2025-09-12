package router

import (
	"YN/config"
	"YN/dao"
	"YN/global"
	"YN/log"
	"YN/model"
	"YN/utils"
	"fmt"
	"net/http"
	"strings"
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
	// 就是我们在RCS三方电梯中，自定义配置的函数名 E1_OutReqTo4F
	deviceID := req.Data.FunctionName[:2]   // E1或E2
	signalType := req.Data.FunctionName[3:] // 信号类型(如OutReqTo4F)

	// 机械臂的光栅独占区异步处理
	if isTrue := rasterExclusiveAreaProcess(req); isTrue {
		return nil
	}

	// 安全检查
	if err := safetyCheck(req, deviceID); err != nil {
		if err.Error() == "req.Data.Value == false" {
			return nil
		}

		return err
	}

	// 业务处理：新增电梯任务
	if err := createElevatorTaskProcess(deviceID, signalType); err != nil {
		return err
	}

	// 业务处理：起始楼层的电梯关门、请求至目标楼层
	if err := elevatorCloseDoorAndToTargetFloorProcess(deviceID, signalType); err != nil {
		return err
	}

	return nil
}

func rasterExclusiveAreaProcess(req fieldFunctionDTO) bool {
	if req.Data.FunctionName == "5M01" { // 站点编号
		global.RasterExclusiveAreaChan1 <- req.Data.Value
		return true
	}

	if req.Data.FunctionName == "5M02" { // 站点编号
		global.RasterExclusiveAreaChan2 <- req.Data.Value
		return true
	}

	return false
}

func safetyCheck(req fieldFunctionDTO, deviceID string) error {
	// RCS的bug：有时候会发Value = false的值，这种请求是不要的
	if req.Data.Value == false {
		return fmt.Errorf("req.Data.Value == false")
	}

	// 非AGV模式 或 电梯故障，不能写入
	if bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])}); bits[4] == 0 || bits[5] == 1 {
		return fmt.Errorf("device %v is manual mode or breakdown! Can't operator Elevator. bits: %v", deviceID, bits)
	}

	// 判断电梯是否在线。读取成功证明在线
	if _, err := global.ClientList[deviceID].ReadDiscreteInputs(config.Config.ReadStartAddr, config.Config.ReadEndAddr); err != nil {
		return fmt.Errorf("elevator %v disconnect", deviceID)
	}

	return nil
}

func createElevatorTaskProcess(deviceID string, signalType string) error {
	if strings.Contains(signalType, "OutReqTo") {
		// 提取起始楼层信息
		startFloorStr := signalType[len(signalType)-2]
		if startFloorStr < '0' || startFloorStr > '9' {
			return fmt.Errorf("startFloorStr = %v, not number", startFloorStr)
		}
		startFloor := float64(startFloorStr - '0')

		// 初始化电梯任务
		elevatorTask := model.ElevatorTask{
			ElevatorID:                      deviceID,
			TaskID:                          deviceID + "_" + utils.GenTaskID(),
			IsProcessToOtherFloorReq:        false,
			IsProcessStartFloorCloseDoorReq: false,
			StartFloor:                      startFloor,
			TargetFloor:                     0,
			TaskStatus:                      global.ElevatorTaskStatus_ToStartFloor,
			TaskType:                        "",
			ReqStatus:                       "",
			StartTime:                       time.Now(),
			UpdateTime:                      time.Now(),
			EndTime:                         time.Time{},
		}
		if err := dao.AddElevatorTaskPool(elevatorTask); err != nil {
			return err
		}
	}

	return nil
}

func elevatorCloseDoorAndToTargetFloorProcess(deviceID string, signalType string) error {
	// () 查询任务，做状态转变
	elevatorTask, exist := dao.GetElevatorTask(deviceID)
	if exist == false {
		return nil
	}

	// () AGV进入电梯后，起始楼层关门处理
	if strings.Contains(signalType, "InReqCloseDoor") && elevatorTask.TaskStatus == global.ElevatorTaskStatus_StartFloorOpenDoorFinish {
		// 写入电梯常关门
		if err := utils.WriteElevatorCoils(deviceID, model.ElevatorSignalCoil{CloseDoor4: 1}); err != nil {
			return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
		}

		// 更新任务状态
		elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_StartFloorCloseDoor, IsProcessStartFloorCloseDoorReq: true}
		if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
			return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
		}
	}

	// () 跨楼层处理
	if strings.Contains(signalType, "InReqTo") {
		i := 0

		for {
			time.Sleep(1 * time.Second)

			i = i + 1
			if i == 120 {
				return nil
			}

			elevatorTaskObj, _ := dao.GetElevatorTask(deviceID)
			// 幂等处理
			if elevatorTask.IsProcessToOtherFloorReq == true {
				return nil
			}
			if elevatorTask.IsProcessStartFloorCloseDoorReq == false {
				return nil
			}

			if elevatorTaskObj.TaskStatus != global.ElevatorTaskStatus_StartFloorCloseDoorFinish {
				continue
			}

			// 提取目标楼层信息
			targetFloorStr := signalType[len(signalType)-2]
			if targetFloorStr < '0' || targetFloorStr > '9' {
				return fmt.Errorf("targetFloorStr = %v, not number", targetFloorStr)
			}
			targetFloor := float64(targetFloorStr - '0')

			// 判断上楼还是下楼
			taskType := global.ElevatorTaskType_Up
			if elevatorTask.StartFloor-targetFloor > 0 {
				taskType = global.ElevatorTaskType_Down
			}

			// 写入电梯去目标楼层
			var signal model.ElevatorSignalCoil
			if targetFloor == 4 {
				signal = model.ElevatorSignalCoil{ReqTo4F1: 1} // 写入电梯去4楼
			}
			if targetFloor == 5 {
				signal = model.ElevatorSignalCoil{ReqTo5F2: 1} // 写入电梯去5楼
			}
			if err := utils.WriteElevatorCoils(deviceID, signal); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
			}

			// 更新任务状态
			elevatorTaskTmp1 := model.ElevatorTask{
				IsProcessToOtherFloorReq: true,
				TargetFloor:              targetFloor,
				TaskType:                 taskType,
				TaskStatus:               global.ElevatorTaskStatus_ToTargetFloor,
			}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTaskTmp1, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}
		}
	}

	// () AGV出电梯后，目标楼层关门处理
	if strings.Contains(signalType, "OutReqCloseDoor") && elevatorTask.TaskStatus == global.ElevatorTaskStatus_TargetFloorOpenDoorFinish {
		// 判断电梯任务池有无任务
		if len(global.ElevatorTaskPool[deviceID]) == 0 { // 没任务
			// 写入电梯常关门
			if err := utils.WriteElevatorCoils(deviceID, model.ElevatorSignalCoil{CloseDoor4: 1}); err != nil {
				return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
			}

			// 更新任务状态
			elevatorTask = model.ElevatorTask{TaskStatus: global.ElevatorTaskStatus_TargetFloorCloseDoor}
			if err := dao.UpdateElevatorTask(deviceID, elevatorTask, false); err != nil {
				return fmt.Errorf("elevator task not found for elevatorID: %v, new-ElevatorTask: %+v. old-ElevatorTaskList: %+v", deviceID, elevatorTask, global.ElevatorTaskList[deviceID])
			}
		} else { // 电梯任务池有任务。则本单任务完成
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
