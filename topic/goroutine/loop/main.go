package main

import "sync"

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(x int) {
			sendRpc(x)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func sendRpc(i int) {
	println(i)
}