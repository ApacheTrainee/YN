package model

type ElevatorSignalCoil struct {
	ReqTo4F1   int // 请求电梯到4楼
	ReqTo5F2   int // 请求电梯到5楼
	OpenDoor3  int // 电梯常开门
	CloseDoor4 int // 电梯常关门
}
