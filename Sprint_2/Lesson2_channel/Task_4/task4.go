package main

import (
	"fmt"
	"sync"
)

// Реализуй паттерны Fan In, Fan Out | Pipeline

// Pipeline
func generator(numbers ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _, number := range numbers {
			out <- number
		}
	}()

	return out
}

// Здесь несколько worker'ов читают один канал — Fan-out
func worker(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for number := range in {
			out <- number * number
		}
	}()

	return out
}

// Fan-in, объединяет несколько каналов в один
func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup

	out := make(chan int)

	wg.Add(len(channels))

	for _, channel := range channels {
		go func() {
			defer wg.Done()

			for value := range channel {
				out <- value
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// pipeline
	numbers := generator(1, 2, 3, 4, 5, 6)

	// fan-in
	worker1 := worker(numbers)
	worker2 := worker(numbers)
	worker3 := worker(numbers)

	// fan-out
	results := merge(worker1, worker2, worker3)

	for result := range results {
		fmt.Println(result)
	}
}
