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
// Задача 10. Параллельная сумма и эксперименты с планировщиком.
//
// 1. Реализуй ParallelSum(nums []int, workers int) int:
//    - срез делится на `workers` примерно равных кусков;
//    - каждый кусок суммируется в своей горутине;
//    - частичные суммы собираются БЕЗ гонок (канал или WaitGroup + Mutex);
//    - workers <= 0 -> использовать runtime.NumCPU();
//    - workers > len(nums) -> не создавать пустых горутин;
//    - результат должен совпадать с последовательной SequentialSum.
//
// 2. Добавь go.mod и файл task10_test.go:
//    - табличный тест (пустой срез, 1 элемент, workers > len, нечётное деление);
//    - бенчмарки BenchmarkSum/seq, BenchmarkSum/par-1, par-2, par-4, par-8
//      на срезе из 10_000_000 элементов. Используй b.Loop() и b.ReportAllocs().
//
// 3. Эксперименты (результаты запиши комментарием в конце файла):
//    - запусти бенчмарк с GOMAXPROCS=1 и с GOMAXPROCS=NumCPU
//      (go test -bench . -cpu 1,2,4,8). Почему при GOMAXPROCS=1 параллельная
//      версия не быстрее (или медленнее) последовательной?
//    - при каком размере среза накладные расходы на горутины перестают
//      окупаться? Найди порог экспериментально и добавь его в ParallelSum
//      как константу minChunk: куски меньше неё не параллелить.
//    - что произойдёт, если в каждой горутине вызывать runtime.Gosched()
//      на каждой итерации? Замерь и объясни.
//
// Проверь: go test -race ./... и go vet ./...

const minChunk = 0 // подбери экспериментально

func SequentialSum(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

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
