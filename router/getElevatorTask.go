package router

import (
	"YN/global"
	"YN/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetElevatorTask(c *gin.Context) {
	req := struct {
		DeviceID string `json:"device_id"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WebLogger.Errorf("SimulateElevatorRead ShouldBindJSON err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "err: Invalid request format",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": global.ElevatorTaskList[req.DeviceID],
	})
}
