package router

import (
	"YN/log"
	"YN/model"
	"YN/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WriteToElevatorDTO struct {
	ReqCode  string `json:"req_code"`
	DeviceID string `json:"device_id"`
	Write1   int    `json:"write_1"`
	Write2   int    `json:"write_2"`
	Write3   int    `json:"write_3"`
	Write4   int    `json:"write_4"`
}

func WriteToElevator(c *gin.Context) {
	var req WriteToElevatorDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WebLogger.Errorf("WriteToElevator ShouldBindJSON err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "err: Invalid request format",
		})
		return
	}
	log.WebLogger.Infof("\n%v \n%+v", c.Request.URL.Path, req.ReqCode)

	if err := writeToElevatorService(req); err != nil {
		log.WebLogger.Errorf("readElevatorStatusService err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("err: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func writeToElevatorService(req WriteToElevatorDTO) error {
	signalCoil := model.ElevatorSignalCoil{
		Write1: req.Write1,
		Write2: req.Write2,
		Write3: req.Write3,
		Write4: req.Write4,
	}
	if err := utils.WriteElevatorCoils(req.DeviceID, signalCoil); err != nil {
		log.Logger.Infof("deviceid = %v reset writeData error: %v", req.DeviceID, err)
		return err
	}

	return nil
}
