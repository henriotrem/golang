package main

import "fmt"

type greeting string

func (g *greeting) Greet() {
	fmt.Println("Bonjour terrien!")
}

var Greeter greeting
