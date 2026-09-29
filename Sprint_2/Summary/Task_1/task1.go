package main

import (
	"fmt"
	"sync"
)

type Job struct {
	ID    int
	Value int
}

type Result struct {
	JobID int
	Value int
}

func WorkerPool(
	workers int,
	jobs []Job,
) []Result {
	var wg sync.WaitGroup
	results := make(chan Result, workers)

	for i := range workers {
		wg.Add(1)
		go worker(jobs[i], results, &wg)
	}

	wg.Wait()
	close(results)

	var result []Result
	for signal := range results {
		result = append(result, signal)
	}
	return result
}

func worker(job Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	results <- Result{
		JobID: job.ID,
		Value: job.Value * job.Value,
	}
}

func main() {
	result := WorkerPool(3, []Job{
		{
			ID:    3,
			Value: 3,
		}, {
			ID:    4,
			Value: 4,
		}, {
			ID:    5,
			Value: 5,
		},
	})

	fmt.Println(result)
}
