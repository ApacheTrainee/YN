package router

import (
	"YN/global"
	"YN/log"
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
	var toFloor float64
	if req.Write1 == 1 {
		toFloor = 4
	}
	if req.Write2 == 1 {
		toFloor = 5
	}

	var doorStatus string
	if req.Write3 == 1 {
		doorStatus = global.OpenDoor
	}
	if req.Write4 == 1 {
		doorStatus = global.CloseDoor
	}

	if err := utils.WriteElevatorCoils(req.DeviceID, toFloor, doorStatus); err != nil {
		log.Logger.Infof("deviceid = %v reset writeData error: %v", req.DeviceID, err)
		return err
	}

	return nil
}
