package service

import (
	"YN/global"
	"YN/log"
	"time"
)

func ElevatorTaskTimeOutProcess() {
	for {
		if global.ElevatorTask.Status != "" {
			if time.Since(global.ElevatorTask.StartTime) > 5*time.Minute {
				log.Logger.Errorf("电梯任务超时，任务状态：%+v", global.ElevatorTask)

				global.ElevatorTask.Status = ""
			}
		}

		time.Sleep(3 * time.Second)
	}
}
