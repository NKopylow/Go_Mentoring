package main

import (
	"fmt"
	"sync"
	"time"
)

// Нужно создать 3 worker'а.

// Каждый worker получает job из канала:

// Каждая задача:

// time.Sleep(500 * time.Millisecond)

// При этом одновременно работает максимум 3 задачи.

// Тренируем: worker pool + jobs channel + results channel + WaitGroup.

func worker(jobsChannel chan int, resultChannel chan int) {
	for job := range jobsChannel {
		time.Sleep(500 * time.Millisecond)
		resultChannel <- job * 2
	}
	time.Sleep(500 * time.Millisecond)
}

func workerPool(jobs []int) []int {
	var wg sync.WaitGroup
	var result []int
	jobChannel := make(chan int)

	go func() {
		for _, job := range jobs {
			jobChannel <- job
		}
		close(jobChannel)
	}()

	resultChannel := make(chan int)

	for range 3 {
		wg.Go(func() {
			worker(jobChannel, resultChannel)
		})
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	for value := range resultChannel {
		result = append(result, value)
	}

	return result
}

func main() {
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(workerPool(jobs))
}
