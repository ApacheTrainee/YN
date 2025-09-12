package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type configFrontFromFile struct {
	AGVPoints  map[string][]string `mapstructure:"agv_points"`
	AGVNumbers []string            `mapstructure:"agv_numbers"`
	MapConfig  map[string]string   `mapstructure:"map_config"`
}

// {AGVPoints:map[4:[4-1 4-2 4-3 4-4 4-5 4-6 4-7 4-8 4-9] 5:[5-1 5-2 5-3 5-4 5-5 5-6 5-7 5-8]] AGVNumbers:[AGV-01 AGV-02 AGV-03 AGV-04] MapConfig:map[map_id_4_floor:四楼1 map_id_5_floor:五楼1]}
func ReadConfigFront(c *gin.Context) {
	v := viper.New()
	v.SetConfigFile("config/config_front.yaml")
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
