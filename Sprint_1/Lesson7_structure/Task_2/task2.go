package main

import (
	"fmt"
	"unsafe"
)

// Создайте две версии структуры с разным порядком полей. С помощью unsafe.Sizeof измерьте их размер и объясните, почему один короче.

type StatsBad struct {
	A bool
	B int64
	C int32
	D bool
}

type StatsGood struct {
	B int64
	C int32
	A bool
	D bool
}

func main() {
	fmt.Println("Bad :", unsafe.Sizeof(StatsBad{}))  // 24
	fmt.Println("Good:", unsafe.Sizeof(StatsGood{})) // 16
}
