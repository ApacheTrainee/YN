package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GenTaskID() string {
	now := time.Now()
	today := now.Format("20060102")

	rand.Seed(now.UnixNano())
	postfix := strconv.Itoa(rand.Intn(1000000))

	return today + postfix
}
