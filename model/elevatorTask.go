package model

import "time"

type ElevatorTask struct {
	Status      string    `json:"status"`
	StartFloor  float64   `json:"startFloor"`
	TargetFloor float64   `json:"targetFloor"`
	StartTime   time.Time `json:"startTime"`
}
