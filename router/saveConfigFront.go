package router

import (
	"YN/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type SaveConfigFrontDTO struct {
	AGVPoints  map[string][]string `json:"AGVPoints"`
	AGVNumbers []string            `json:"AGVNumbers"`
	MapConfig  map[string]string   `json:"MapConfig"`
}

func SaveConfigFront(c *gin.Context) {
	var req SaveConfigFrontDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WebLogger.Errorf("SaveConfigFront ShouldBindJSON err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "err: Invalid request format",
		})
		return
	}
	log.WebLogger.Infof("\n%v \n", c.Request.URL.Path)

	if err := saveConfigFrontService(req); err != nil {
		log.WebLogger.Errorf("saveConfigFrontService err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": fmt.Sprintf("err: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func saveConfigFrontService(req SaveConfigFrontDTO) error {
	v := viper.New()
	v.SetConfigFile("config/config_front.yaml")
	v.Set("agv_points", req.AGVPoints)
	v.Set("agv_numbers", req.AGVNumbers)
	v.Set("map_config", req.MapConfig)
	if err := v.WriteConfig(); err != nil {
		return err
	}

	return nil
}
