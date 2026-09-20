package main

import (
	"fmt"
	"sync"
)

func writer(channel chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	channel <- 500
}

func main() {
	channel := make(chan int)
	var result int
	var wg sync.WaitGroup

	wg.Add(1)
	go writer(channel, &wg)

	go func() {
		defer close(channel)
		result = <-channel
	}()

	wg.Wait()

	fmt.Print("Value is ", result)
}
