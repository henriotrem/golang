package raft

import (
	"log"
	"sort"
)

// Debugging
const Debug = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func CopyAndSort(source []int) []int {
	destination := make([]int, len(source))
	copy(destination, source)
	sort.Ints(destination)

	return destination
}
