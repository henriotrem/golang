package main

import "fmt"

type entier int
type Entier = int

func main() {

	var x int = 12
	var y Entier = 10
	var z = sum(x, y)

	fmt.Printf("%T\n%T\n%T\n%v", x, y, z, z)
}

func sum(x, y int) entier {
	return entier(x + y)
}
