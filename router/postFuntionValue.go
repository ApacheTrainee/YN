package router

import (
	"YN/config"
	"YN/global"
	"YN/log"
	"YN/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
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

	// () 电梯去起始楼层
	if strings.Contains(signalType, "OutReqTo") {
		// 提取起始楼层信息
		startFloorStr := signalType[len(signalType)-2]
		if startFloorStr < '0' || startFloorStr > '9' {
			return fmt.Errorf("startFloorStr = %v, not number", startFloorStr)
		}
		startFloor := float64(startFloorStr - '0')

		global.ElevatorTask.Status = global.ElevatorTaskStatus_ToStartFloor
		global.ElevatorTask.StartFloor = startFloor
		global.ElevatorTask.StartTime = time.Now()

		global.StartFloorProcessChan <- startFloor
		for result := range global.StartFloorProcessChanResult {
			if result == "ok" {
				return nil
			} else {
				return fmt.Errorf("startFloor = %v, error: %v", startFloor, result)
			}
		}
	}

	// () AGV进入电梯后，起始楼层关门处理
	// 起始楼层关门，和去目标楼层时同时发的，直接在目标楼层那里处理，关门和去目标楼层写入了
	if strings.Contains(signalType, "InReqCloseDoor") {
		//// 写入电梯常关门
		//if err := utils.WriteElevatorCoils(deviceID, 0, global.CloseDoor); err != nil {
		//	return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
		//}
	}

	// () 跨楼层处理。会重复发
	if strings.Contains(signalType, "InReqTo") {
		// 提取起始楼层信息
		targetFloorStr := signalType[len(signalType)-2]
		if targetFloorStr < '0' || targetFloorStr > '9' {
			return fmt.Errorf("targetFloorStr = %v, not number", targetFloorStr)
		}
		targetFloor := float64(targetFloorStr - '0')

		global.ElevatorTask.TargetFloor = targetFloor
		global.ElevatorTask.Status = global.ElevatorTaskStatus_ToTargetFloor
		global.ElevatorTask.StartTime = time.Now()

		// 写入电梯去目标楼层
		if err := utils.WriteElevatorCoils(deviceID, targetFloor, global.CloseDoor); err != nil {
			return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
		}
	}

	// ()  AGV出了电梯后，目标楼层关门处理
	if strings.Contains(signalType, "OutReqCloseDoor") {
		// 写入电梯常关门
		if err := utils.WriteElevatorCoils(deviceID, 0, global.CloseDoor); err != nil {
			return fmt.Errorf("deviceid = %v reset writeData error: %v", deviceID, err)
		}
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

	// 判断电梯是否在线。读取成功证明在线
	if _, err := global.ClientList[deviceID].ReadDiscreteInputs(config.Config.ReadStartAddr, config.Config.ReadEndAddr); err != nil {
		return fmt.Errorf("elevator %v disconnect", deviceID)
	}

	// 校验：读取电梯的值是否正确
	bits := utils.BytesToBits([]byte{byte(global.ElevatorStatus[deviceID])})
	if (bits[0] == 0 && bits[1] == 0 && bits[2] == 0 && bits[3] == 0) || (bits[0] == 1 && bits[1] == 1 && bits[2] == 1 && bits[3] == 1) || bits[4] == 0 || bits[5] == 1 {
		return fmt.Errorf("elevator err. bits: %v", bits)
	}

	return nil
}
