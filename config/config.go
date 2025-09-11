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
		panic("read config file err: " + err.Error())
	}

	if err := v.Unmarshal(&Config); err != nil {
		panic("Unmarshal config file err: " + err.Error())
	}
}
