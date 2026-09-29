package main

import (
	"context"
	"fmt"
	"sync"
	"time"
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
	ctx context.Context,
) []Result {
	var wg sync.WaitGroup

	jobsCh := make(chan Job)
	results := make(chan Result)

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go worker(ctx, jobsCh, results, &wg)
	}

	go func() {
		defer close(jobsCh)

		for _, job := range jobs {
			select {
			case jobsCh <- job:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var result []Result

	for res := range results {
		result = append(result, res)
	}

	return result
}

func worker(
	ctx context.Context,
	jobs <-chan Job,
	results chan<- Result,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			select {
			case results <- Result{
				JobID: job.ID,
				Value: job.Value * job.Value,
			}:
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Millisecond,
	)
	defer cancel()
	slice := []Job{
		{
			ID:    3,
			Value: 3,
		},
		{
			ID:    4,
			Value: 4,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
		{
			ID:    5,
			Value: 5,
		},
	}
	result := WorkerPool(len(slice), slice, ctx)

	fmt.Println(result)
}
