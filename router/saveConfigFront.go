package router

import (
	"YN/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path/filepath"
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
	// 获取可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		log.WebLogger.Panicf("Failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath) // 获取可执行文件所在目录

	// 基于可执行文件目录构造配置文件的绝对路径
	configFilePath := filepath.Join(exeDir, "config", "config.yaml")

	v := viper.New()
	v.SetConfigFile(configFilePath)
	v.Set("agv_points", req.AGVPoints)
	v.Set("agv_numbers", req.AGVNumbers)
	v.Set("map_config", req.MapConfig)
	if err := v.WriteConfig(); err != nil {
		return err
	}

	return nil
}
