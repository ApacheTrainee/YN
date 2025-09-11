package router

import (
	"YN/global"
	"YN/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SimulateElevatorRead(c *gin.Context) {
	req := struct {
		CoilValue []byte `json:"coilValue"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WebLogger.Errorf("SimulateElevatorRead ShouldBindJSON err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "err: Invalid request format",
		})
		return
	}
	log.WebLogger.Infof("\n%v \n%+v", c.Request.URL.Path, req.CoilValue)
	log.Logger.Infof("\n%v \n%+v", c.Request.URL.Path, req.CoilValue)

	global.ElevatorCoilSimulation = req.CoilValue

	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
