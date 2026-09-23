package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	var data [][]byte

	for i := 0; i < 100_000; i++ {
		data = append(data, make([]byte, 1024))
	}

	start := time.Now()

	runtime.GC()

	elapsed := time.Since(start)

	fmt.Println("GC took:", elapsed)

	for i := 0; i < 5; i++ {
		start := time.Now()

		runtime.GC()

		fmt.Println(time.Since(start))
	}
}
