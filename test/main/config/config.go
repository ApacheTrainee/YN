package config

import (
	"github.com/spf13/viper"
)

var Config ConfigFromFile

type ConfigFromFile struct {
	DeviceList    []Device `mapstructure:"deviceList"`
	ReadStartAddr uint16   `mapstructure:"readStartAddr"`
	ReadEndAddr   uint16   `mapstructure:"readEndAddr"`
	RcsIP         string   `mapstructure:"rcsIP"`
	RcsPort       string   `mapstructure:"rcsPort"`
	RunMode       string   `mapstructure:"runMode"`
}

type Device struct {
	Id      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("config/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic("read config file err: " + err.Error())
	}

	if err := v.Unmarshal(&Config); err != nil {
		panic("Unmarshal config file err: " + err.Error())
	}
}
