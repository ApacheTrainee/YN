package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var Config ConfigFromFile

type ConfigFromFile struct {
	DeviceList    []Device `mapstructure:"deviceList"`
	ReadStartAddr uint16   `mapstructure:"readStartAddr"`
	ReadEndAddr   uint16   `mapstructure:"readEndAddr"`
	RcsIP         string   `mapstructure:"rcsIP"`
	RcsPort       string   `mapstructure:"rcsPort"`
	AllPartIP     string   `mapstructure:"allPartIP"`
	AllPartPort   string   `mapstructure:"allPartPort"`
	RunMode       string   `mapstructure:"runMode"`
}

type Device struct {
	Id      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func InitConfig() {
	currentDir, errPath := os.Getwd() // 获取当前目录
	if errPath != nil {
		panic("Failed to obtain the absolute path of the current directory" + errPath.Error())
	}

	configFilePath := filepath.Join(currentDir, "config", "config.yaml") // 定义日志输入路径

	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		log.Panicf("config file not exist: %s", configFilePath)
	}

	v := viper.New()
	v.SetConfigFile(configFilePath)
	if err := v.ReadInConfig(); err != nil {
		panic("read config file err: " + err.Error())
	}

	if err := v.Unmarshal(&Config); err != nil {
		panic("Unmarshal config file err: " + err.Error())
	}
}
