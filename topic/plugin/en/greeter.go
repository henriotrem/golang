package main

import "fmt"

type greeting string

func (g *greeting) Greet() {
	fmt.Println("Hello world!")
}

var Greeter greeting
