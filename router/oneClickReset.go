package router

import (
	"YN/global"
	"YN/log"
	"YN/model"
	"YN/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type oneClickReset struct {
	ReqCode                string `json:"req_code"`
	DeviceID               string `json:"device_id"`
	OnlyDeletePlcWrite     bool   `json:"only_delete_plc_write"`
	OnlyDeleteElevatorTask bool   `json:"only_delete_elevator_task"`
}

func OneClickReset(c *gin.Context) {
	var req oneClickReset
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WebLogger.Errorf("OneClickReset ShouldBindJSON err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "err: Invalid request format",
		})
		return
	}
	log.WebLogger.Infof("\n%v \n%+v", c.Request.URL.Path, req.ReqCode)

	if err := OneClickResetService(req); err != nil {
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

func OneClickResetService(req oneClickReset) error {
	if req.OnlyDeleteElevatorTask == true {
		delete(global.ElevatorTaskList, req.DeviceID)
		return nil
	}
	if req.OnlyDeletePlcWrite == true {
		// 清空电梯写入
		if err := utils.WriteElevatorCoils(req.DeviceID, model.ElevatorSignalCoil{}); err != nil {
			log.Logger.Infof("deviceid = %v reset writeData error: %v", req.DeviceID, err)
			return err
		}
		return nil
	}

	delete(global.ElevatorTaskList, req.DeviceID)

	// 清空电梯写入
	if err := utils.WriteElevatorCoils(req.DeviceID, model.ElevatorSignalCoil{}); err != nil {
		log.Logger.Infof("deviceid = %v reset writeData error: %v", req.DeviceID, err)
		return err
	}

	return nil
}
