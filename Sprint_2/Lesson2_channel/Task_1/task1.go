package main

import (
	"fmt"
	"sync"
)

// Реализуй программу,
// где две горутины отправляют числа в канал, а основная функция читает и печатает их.

func task1(channel chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	channel <- 100
}

func task2(channel chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	channel <- 500
}

func main() {
	channel := make(chan int)
	var wg sync.WaitGroup

	wg.Add(2)

	go task1(channel, &wg)
	go task2(channel, &wg)

	go func() {
		wg.Wait()
		close(channel)
	}()

	// select {
	// case value := <-channel:
	// 	fmt.Printf("Value is: ", value)
	// }

	for value := range channel {
		fmt.Printf("Value is: ", value)
	}
}
