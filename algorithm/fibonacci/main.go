package main

import "fmt"

func getNthFib(number int) int {
	if number == 1 {
		return 0
	}
	f1, f2 := 0, 1
	for i := 1; i < number-1; i++ {
		f1, f2 = f2, f1+f2
	}
	return f2
}

func main() {
	for i := 1; i < 10; i++ {
		fmt.Println(getNthFib(i))
	}
}
