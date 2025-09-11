package utils

import (
	"YN/global"
	"YN/log"
	"YN/model"
	"fmt"
	"strconv"
)

// 写入电梯的coil信号
func WriteElevatorCoils(deviceId string, signal model.ElevatorSignalCoil) error {
	coils := make([]int, 4)
	coils[0] = signal.ReqTo4F1
	coils[1] = signal.ReqTo5F2
	coils[2] = signal.OpenDoor3
	coils[3] = signal.CloseDoor4

	log.WebLogger.Infof("deviceId: %v, preparing to write coils: %v", deviceId, coils)

	// 将int切片，转换为只有一个byte值的切片
	coilsBytes, err := intSliceToOneByteSlice(coils)
	if err != nil {
		return fmt.Errorf("boolSliceToByteSlice err: %v", err)
	}

	// 写入电梯信号
	client, ok := global.ClientList[deviceId]
	if !ok {
		return fmt.Errorf("device %s not connected", deviceId)
	}

	_, err = client.WriteMultipleCoils(0, 4, coilsBytes)
	if err != nil {
		return fmt.Errorf("failed to write coils to device = %v err: %v", deviceId, err)
	}
	log.WebLogger.Infof("write to device %v Successfully\n", deviceId)

	return nil
}

func intSliceToOneByteSlice(coils []int) ([]byte, error) {
	// 反转数组，方式一
	for i, j := 0, len(coils)-1; i < j; i, j = i+1, j-1 {
		tmp := coils[j]
		coils[j] = coils[i]
		coils[i] = tmp
		//coils[i], coils[j] = coils[j], coils[i] // 交换首尾元素
	}

	//// 反转数组，方式二
	//coilsTmp := make([]int, len(coils))
	//for i := 0; i < len(coils); i++ {
	//	coilsTmp[i] = coils[len(coils)-1-i] // 反向填充新切片
	//}

	var binStr string
	for _, value := range coils {
		binStr = binStr + fmt.Sprintf("%v", value)
	}

	num, err := strconv.ParseInt(binStr, 2, 8)
	if err != nil {
		return nil, fmt.Errorf("strconv.ParseInt err: %v", err)
	}
	result := []byte{byte(num)}

	return result, nil
}
