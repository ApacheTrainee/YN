package model

type ElevatorSignalCoil struct {
	Write1 int // 请求电梯到4楼
	Write2 int // 请求电梯到5楼
	Write3 int // 电梯常开门
	Write4 int // 电梯常关门
}
