package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
)

// Тема: горутины и планировщик.
//
// Задача 10. Параллельная сумма.
//
// Реализуй ParallelSum. Добавь go.mod + task10_test.go:
//   - табличный тест: пустой срез, 1 элемент, workers > len, нечётное деление;
//   - бенчмарки seq / par-1 / par-4 / par-8 (b.Loop, b.ReportAllocs) на 10_000_000 элементов.
//
// Эксперимент (результат — комментарием в конце файла):
//   go test -bench . -cpu 1,4 — почему при GOMAXPROCS=1 параллельная версия не быстрее?
//
// Проверь: go test -race ./...

const minChunk = 0 // опционально: не параллелить, если кусок меньше

// SequentialSum суммирует nums в одном потоке. Эталон для проверки ParallelSum.
func SequentialSum(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

// ParallelSum делит nums на workers кусков и суммирует каждый в своей горутине.
//
// Поведение:
//   - workers <= 0  → runtime.NumCPU()
//   - workers > len(nums) → не создавать пустых горутин (уменьшить workers)
//   - частичные суммы собрать без гонок (канал / WaitGroup+Mutex / atomic)
//   - результат == SequentialSum(nums)
func ParallelSum(nums []int, workers int) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Catching panic")
		}
	}()
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	if workers > len(nums) {
		workers = len(nums)
	}
	var sum atomic.Int64
	wg := sync.WaitGroup{}
	chuncks := int(math.Ceil(float64(len(nums)) / float64(workers)))
	for i := 0; i < workers; i++ {
		wg.Go(func() {
			end := (i + 1) * chuncks
			if end > len(nums) {
				end = len(nums)
			}
			chunk := nums[i*chuncks : end]
			fmt.Println("i*chuncks:", i*chuncks, "(i+1)*chuncks: ", (i+1)*chuncks)
			chunkSum := 0
			for _, n := range chunk {
				chunkSum += n
			}
			sum.Add(int64(chunkSum))
		})
	}
	// panic("not implemented")
	wg.Wait()
	return int(sum.Load())
}

func main() {
	nums := make([]int, 1_000_000)
	for i := range nums {
		nums[i] = i
	}

	fmt.Println("NumCPU:", runtime.NumCPU(), "GOMAXPROCS:", runtime.GOMAXPROCS(0))
	fmt.Println("seq:", SequentialSum(nums))
	fmt.Println("par:", ParallelSum(nums, 7))
}

// Результаты экспериментов:
//
