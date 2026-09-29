package main

import (
	"context"
	"fmt"
	"sync"
)

func worker(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job := <-jobs:
			fmt.Println("Job has arrived, her id is: ", id, "job result: ", job)
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	jobs := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := range 3 {
		wg.Add(1)
		go worker(ctx, i, jobs, &wg)
		select {
		case jobs <- i * 2:
		case <-ctx.Done():
			break
		}
	}
	cancel()
	wg.Wait()
}
