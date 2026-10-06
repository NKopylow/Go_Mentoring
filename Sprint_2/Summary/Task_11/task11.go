package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Тема: каналы (pipeline), context.
//
// Задача 11. Конвейер с отменой.
//
// Реализуй стадии ниже. Общие правила для ВСЕХ стадий:
//   - чтение/запись в каналы — через select с ctx.Done() (не блокироваться навсегда при отмене);
//   - при завершении закрывать свой выходной канал;
//   - не оставлять утекающих горутин.
//
// В main: Generate → Filter(чётные) → FanOut(4, slowSquare) → Take(5), вывести результат.
// Проверь с -race.

// Generate запускает горутину, которая по очереди отправляет nums в канал и закрывает его.
// Если ctx отменён до/во время отправки — прекращает работу и закрывает канал.
func Generate(ctx context.Context, nums ...int) <-chan int {
	jobs := make(chan int)
	var wg sync.WaitGroup

	select {
	case <-ctx.Done():
		close(jobs)
		return jobs
	default:
		wg.Go(func() {
			for _, value := range nums {
				jobs <- value
			}
		})

		go func() {
			wg.Wait()
			close(jobs)
		}()
		return jobs
	}
}

// Map читает значения из in, применяет fn к каждому и пишет результат в выходной канал.
// Завершается, когда in закрыт или ctx отменён.
func Map[T, R any](ctx context.Context, in <-chan T, fn func(T) R) <-chan R {
	panic("not implemented")
}

// Filter читает значения из in и пропускает дальше только те, для которых pred == true.
// Завершается, когда in закрыт или ctx отменён.
func Filter[T any](ctx context.Context, in <-chan T, pred func(T) bool) <-chan T {
	panic("not implemented")
}

// FanOut запускает workers горутин: каждая читает из in, применяет fn, пишет в общий выход.
// Порядок результатов НЕ важен. Когда in закрыт и все воркеры закончили — закрыть выход.
func FanOut[T, R any](ctx context.Context, in <-chan T, workers int, fn func(T) R) <-chan R {
	panic("not implemented")
}

// Take пропускает первые n значений из in и закрывает выход.
// После этого продюсеры выше по цепочке должны завершиться (через отмену ctx), а не зависнуть.
func Take[T any](ctx context.Context, in <-chan T, n int) <-chan T {
	panic("not implemented")
}

func slowSquare(n int) int {
	time.Sleep(100 * time.Millisecond)
	return n * n
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	nums := make([]int, 0, 100)
	for i := 1; i <= 100; i++ {
		nums = append(nums, i)
	}

	even := Filter(ctx, Generate(ctx, nums...), func(n int) bool { return n%2 == 0 })
	squares := FanOut(ctx, even, 4, slowSquare)

	for v := range Take(ctx, squares, 5) {
		fmt.Println(v)
	}
}
