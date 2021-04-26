package main

import (
	"fmt"
)

func main() {

	var array = [10]int{4, 2, 8, 1, 3, 0, 7, 9, 5, 6}

	for i := 1; i < len(array); i++ {
		for j := i; j > 0 && array[j-1] > array[j]; j-- {
			array[j-1], array[j] = array[j], array[j-1]
		} 
	}

	fmt.Println(array)
}