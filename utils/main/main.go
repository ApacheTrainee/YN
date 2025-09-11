package main

import (
	"os"
	"path/filepath"
)

type a struct {
	Age string
}

func main() {
	// 定义日志输入路径
	currentDir, errPath := os.Getwd() // 获取当前目录
	if errPath != nil {
		panic("Failed to obtain the absolute path of the current directory" + errPath.Error())
	}

	var logOutputPath string
	logOutputPath = filepath.Join(currentDir, "log", "logfile") // 定义日志输入路径

	if _, err := os.Stat(logOutputPath); os.IsNotExist(err) {
		parentDir := filepath.Dir(currentDir)                      // 获取上一级目录
		logOutputPath = filepath.Join(parentDir, "log", "logfile") // 定义日志输入路径
	}
}
