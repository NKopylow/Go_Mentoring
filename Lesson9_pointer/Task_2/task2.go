package main

import (
	"fmt"
	"unsafe"
)

// Напиши программу, выводящую размер указателя на int, string и пустую структуру.

func main() {
	// TODO: выведи три размера
	k := new(int)
	v := new(string)
	m := new(struct{})
	fmt.Println(unsafe.Sizeof(k))
	fmt.Println(unsafe.Sizeof(v))
	fmt.Println(unsafe.Sizeof(m))
}
