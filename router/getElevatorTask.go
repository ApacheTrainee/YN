package router

import (
	"YN/global"
	"YN/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetElevatorTask(c *gin.Context) {
	getElevatorTaskDTO := struct {
		DeviceID string `json:"device_id"`
	}{}
	if err := c.ShouldBindJSON(&getElevatorTaskDTO); err != nil {
		log.WebLogger.Errorf("GetElevatorTask ShouldBindJSON err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "err: Invalid request format",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": global.ElevatorTask})
}
