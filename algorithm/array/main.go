package main

import (
	"fmt"
	"sort"
)

func main() {
	array := []int{3, 5, -4, 8, 11, 1, -1, 6}
	target := 10

	fmt.Println(TwoNumberSum1(array, target))
	fmt.Println(TwoNumberSum2(array, target))
}

func TwoNumberSum1(array []int, target int) []int {
	exist := make(map[int]bool)
	for _, elem := range array {
		if exist[target-elem] {
			return []int{target - elem, elem}
		}
		exist[elem] = true
	}
	return []int{}
}

func TwoNumberSum2(array []int, target int) []int {
	sort.Ints(array)
	left, right := 0, len(array)-1
	for left < right {
		result := array[left] + array[right]
		if result == target {
			return []int{array[left], array[right]}
		} else if result < target {
			left++
		} else {
			right--
		}
	}
	return []int{}
}
