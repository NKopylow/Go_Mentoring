package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Итоговая задача (все темы спринта по чуть-чуть).
//
// Задача 21. Мини-планировщик: очередь задач + N воркеров.
//
// Это похоже на твой Task_7 (worker pool), но с context и статусами.
//
// Шаги:
//   1. Реализуй Start: запустить workers горутин-воркеров.
//      Каждый воркер в цикле читает задачу из канала jobs и вызывает job.Fn(ctx).
//   2. Реализуй Submit: положить задачу в канал jobs.
//   3. Реализуй Stop: закрыть канал jobs, дождаться воркеров через WaitGroup.
//   4. Запусти main — все 5 задач должны напечатать "done".
//   5. Проверь: go run -race .
//
// Усложнение (по желанию):
//   - если ctx в Submit уже отменён — не принимать задачу;
//   - в каждой задаче проверять <-ctx.Done().

type Job struct {
	Name string
	Fn   func(ctx context.Context) error
}

type Scheduler struct {
	workers int
	jobs    chan Job
	wg      sync.WaitGroup
}

// New создаёт планировщик.
// workers — сколько горутин будут брать задачи из очереди.
// Канал jobs сделай буферизованным (например, на 64).
func New(workers int) *Scheduler {
	panic("not implemented")
}

// Start запускает s.workers горутин.
// Каждая:
//
//	for job := range s.jobs {
//	    _ = job.Fn(ctx)
//	}
//
// Не забудь wg.Add / wg.Done.
func (s *Scheduler) Start(ctx context.Context) {
	panic("not implemented")
}

// Submit кладёт job в очередь (s.jobs <- job).
// Если хочешь поддержать отмену — используй select с ctx.Done().
func (s *Scheduler) Submit(ctx context.Context, job Job) error {
	panic("not implemented")
}

// Stop закрывает канал jobs и ждёт, пока все воркеры закончат (wg.Wait).
// После Stop Submit больше вызывать нельзя.
func (s *Scheduler) Stop() {
	panic("not implemented")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := New(3)
	s.Start(ctx)

	for i := 1; i <= 5; i++ {
		id := i
		_ = s.Submit(ctx, Job{
			Name: fmt.Sprintf("job-%d", id),
			Fn: func(ctx context.Context) error {
				time.Sleep(200 * time.Millisecond)
				fmt.Println("done", id)
				return nil
			},
		})
	}

	s.Stop()
	fmt.Println("all stopped")
}
