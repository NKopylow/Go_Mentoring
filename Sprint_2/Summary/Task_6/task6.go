package main

import (
	"fmt"
	"sync"
	"time"
)

// Есть 10 задач:
// Каждая задача занимает 1 секунду:
// Нужно сделать так, чтобы одновременно выполнялось максимум 5 задач.

func process(id int, channel chan int) {
	fmt.Println("start", id)
	time.Sleep(time.Second * 5)
	fmt.Println("finish", id)
	fmt.Println("-------------------------")
	channel <- id * 2
}

func main() {
	sem := make(chan struct{}, 5)
	var result []int
	valuesChan := make(chan int)
	var wg sync.WaitGroup
	// var mu sync.Mutex
	defer close(sem)

	jobs := []int{
		1, 2, 3, 4, 5,
		6, 7, 8, 9, 10,
	}

	for _, value := range jobs {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() {
				<-sem
			}()
			process(value, valuesChan)
		})
	}

	go func() {
		wg.Wait()
		close(valuesChan)
	}()

	fmt.Println("Waiting")

	for value := range valuesChan {
		result = append(result, value)
	}

	fmt.Println(result)
}
