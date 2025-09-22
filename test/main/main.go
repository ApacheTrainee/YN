package main

import (
	"fmt"
	"slices" // 需要导入 slices 包
)

func main() {
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	slice3 := []int{1, 2, 4}

	fmt.Println("slice1 == slice2:", slices.Equal(slice1, slice2)) // true
	fmt.Println("slice1 == slice3:", slices.Equal(slice1, slice3)) // false
}
