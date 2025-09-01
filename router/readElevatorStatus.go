package router

import (
	"YN/config"
	"YN/global"
	"YN/log"
	"YN/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type readElevatorStatusDTO struct {
	ReqCode  string `json:"req_code"`
	DeviceID string `json:"device_id"`
}

func ReadElevatorStatus(c *gin.Context) {
	var req readElevatorStatusDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WebLogger.Errorf("ReadElevatorStatus ShouldBindJSON err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "err: Invalid request format",
		})
		return
	}
	log.WebLogger.Infof("\n%v \n%+v", c.Request.URL.Path, req.ReqCode)

	inputResult, coilResult, err := readElevatorStatusService(req)
	if err != nil {
		log.WebLogger.Errorf("readElevatorStatusService err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("err: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"req_code":          req.ReqCode,
		"device_id":         req.DeviceID,
		"read_input_result": inputResult,
		"write_coil_result": coilResult,
	})
}

func readElevatorStatusService(req readElevatorStatusDTO) ([]int, []int, error) {
	// plc地址值不开，读取会报错
	inputsResult, err := global.ClientList[req.DeviceID].ReadDiscreteInputs(config.Config.ReadStartAddr, config.Config.ReadEndAddr)
	if err != nil {
		return nil, nil, err
	}

	coilsResults, err := global.ClientList[req.DeviceID].ReadCoils(0, 4)
	if err != nil {
		return nil, nil, err
	}

	inputResult := utils.BytesToBits(inputsResult)
	coilResult := utils.BytesToBits(coilsResults)
	return inputResult, coilResult[0:4], nil
}
