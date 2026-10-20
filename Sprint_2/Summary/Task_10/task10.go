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
// Цель: научиться делить работу между горутинами и собирать результат без гонок.
//
// Шаги:
//   1. Допиши ParallelSum (см. описание функции ниже).
//   2. В main сравни SequentialSum и ParallelSum — числа должны совпасть.
//   3. (по желанию) Добавь go.mod и простой тест TestParallelSum.
//   4. Запусти: go test -race .   и   go test -bench . -cpu 1,4
//      Напиши комментарием в конце: почему при -cpu 1 параллельная версия
//      не быстрее последовательной?

// SequentialSum — обычная сумма в одном потоке. Уже готова, не меняй.
// Нужна, чтобы сверить результат ParallelSum.
func SequentialSum(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

// ParallelSum должна:
//  1. Разделить nums на workers примерно равных кусков.
//  2. Каждый кусок посчитать в своей горутине.
//  3. Сложить частичные суммы (atomic / Mutex / канал — на выбор).
//  4. Дождаться всех горутин (WaitGroup) и вернуть итог.
//
// Краевые случаи:
//   - workers <= 0           → взять runtime.NumCPU()
//   - workers > len(nums)    → уменьшить workers до len(nums)
//   - пустой nums            → вернуть 0
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

// Почему при -cpu 1 параллельная версия не быстрее:
//
