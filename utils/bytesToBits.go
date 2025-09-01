package utils

import "YN/config"

// bytes转为Bits
func BytesToBits(data []byte) []int {
	bitList := make([]int, config.Config.ReadEndAddr) // 创建一个包含所有位的列表

	addrLen := int(config.Config.ReadEndAddr)
	for index, value := range data {
		for i := 0; i < addrLen; i++ {
			bitList[index*addrLen+i] = int((value >> i) & 1) // 提取每个位
		}
	}

	return bitList
}
