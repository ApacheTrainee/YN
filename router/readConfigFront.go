package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type configFrontFromFile struct {
	AGVPoints  map[string][]string `mapstructure:"agv_points"`
	AGVNumbers []string            `mapstructure:"agv_numbers"`
	MapConfig  map[string]string   `mapstructure:"map_config"`
}

// {AGVPoints:map[4:[4-1 4-2 4-3 4-4 4-5 4-6 4-7 4-8 4-9] 5:[5-1 5-2 5-3 5-4 5-5 5-6 5-7 5-8]] AGVNumbers:[AGV-01 AGV-02 AGV-03 AGV-04] MapConfig:map[map_id_4_floor:四楼1 map_id_5_floor:五楼1]}
func ReadConfigFront(c *gin.Context) {
	// 获取可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		log.Panicf("Failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath) // 获取可执行文件所在目录

	// 基于可执行文件目录构造配置文件的绝对路径
	configFilePath := filepath.Join(exeDir, "config", "config.yaml")

	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		log.Panicf("config file not exist: %s", configFilePath)
	}

	v := viper.New()
	v.SetConfigFile(configFilePath)
	if err := v.ReadInConfig(); err != nil {
		panic("read config_front file err: " + err.Error())
	}

	configFront := configFrontFromFile{}
	if err := v.Unmarshal(&configFront); err != nil {
		panic("Unmarshal config_front file err: " + err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"AGVPoints":  configFront.AGVPoints,
		"AGVNumbers": configFront.AGVNumbers,
		"MapConfig":  configFront.MapConfig,
	})
}
