package service

import (
	"YN/config"
	"YN/global"
	"YN/log"
	"YN/utils"
	"fmt"
	"github.com/goburrow/modbus"
	"sync"
	"time"
)

func StartEquipmentMonitor() {
	var wg sync.WaitGroup
	for _, device := range config.Config.DeviceList {
		wg.Add(1)
		go func(dev config.Device) {
			defer wg.Done()

			equipmentStatusMonitor(dev)
		}(device)
	}
	// todo 处理断线检测?

	wg.Wait()
}

// 设备连接、状态逻辑处理循环
func equipmentStatusMonitor(device config.Device) {
	// () 设备连接
	address := fmt.Sprintf("%s:%d", device.Address, device.Port)
	handler := modbus.NewTCPClientHandler(address)
	defer handler.Close()
	handler.Timeout = 3 * time.Second
	handler.SlaveId = 1

	for {
		if err := handler.Connect(); err != nil {
			log.Logger.Errorf("connection to the elevator failed. err: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Logger.Infof("connection to the elevator successfully")

		client := modbus.NewClient(handler)
		global.ClientList[device.Id] = client // 使用设备ID作为键
		break
	}

	// () 设备读取初始化操作
	// result = [32]  2的次方
	result, _ := global.ClientList[device.Id].ReadDiscreteInputs(config.Config.ReadStartAddr, config.Config.ReadEndAddr)
	preResult := result

	// 转为二进制，左往右数	bits = [32] = [0 0 0 0 1 0]
	bits := utils.BytesToBits(preResult)
	log.Logger.Infof("[Device status Change] DeviceID:%v, result = %v. BytesToBits %v", device.Id, result, bits)

	signalValue := int(result[0])
	global.ElevatorStatus[device.Id] = signalValue

	// 循环读取状态，根据是否变化，处理电梯信号
	for {
		// 只读：modsim32 for opto22的input status
		result, err := global.ClientList[device.Id].ReadDiscreteInputs(config.Config.ReadStartAddr, config.Config.ReadEndAddr)
		if err != nil {
			log.Logger.Infof("ReadDiscreteInputs err: %v", err)
			time.Sleep(3 * time.Second)

			for {
				if err = handler.Connect(); err != nil {
					log.Logger.Errorf("connection to the elevator failed. err: %v", err)
					time.Sleep(3 * time.Second)
					continue
				}
				log.Logger.Infof("connection to the elevator successfully")

				client := modbus.NewClient(handler)
				global.ClientList[device.Id] = client // 使用设备ID作为键
				break
			}

			continue
		}

		if result[0] != preResult[0] {
			bits = utils.BytesToBits(result)
			log.Logger.Infof("[Device status Change] DeviceID:%v, result = %v. BytesToBits %v", device.Id, result, bits)

			signalValue = int(result[0])
			global.ElevatorStatus[device.Id] = signalValue

			processElevatorSignal(device, signalValue) // 处理电梯信号
		}
		preResult = result

		time.Sleep(1 * time.Second)
	}
}
