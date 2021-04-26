package main

import (
	"fmt"
	"net"
)

func worker(input, output chan int) {
	for port := range input {
		address := fmt.Sprintf("scanme.nmap.com:%d", port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			output <- 0
			continue
		}
		conn.Close()
		output <- port
	}
}

func main() {
	input := make(chan int, 100)
	output := make(chan int)
	var openports []int

	for i := 0; i < cap(input); i++ {
		go worker(input, output)
	}

	for i := 0; i < 1024; i++ {
		go func(port int) {
			input <- port
		}(i)
	}

	for i := 0; i < 1024; i++ {
		port := <-output
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(input)
	close(output)
	fmt.Println(len(openports), openports)
}
