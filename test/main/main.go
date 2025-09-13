package main

import (
	"YN/model"
	"fmt"
)

type a struct {
	Age string
}

func main() {
	signal := model.ElevatorSignalCoil{}

	coils := make([]int, 4)
	coils[0] = signal.ReqTo4F1
	coils[1] = signal.ReqTo5F2
	coils[2] = signal.OpenDoor3
	coils[3] = signal.CloseDoor4

	if coils[0] == 0 && coils[1] == 0 && coils[2] == 0 && coils[3] == 0 {
		fmt.Printf("=== %v", len(coils))
	}
}
