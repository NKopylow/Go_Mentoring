package main

import (
	"fmt"

	"stack/stack"
)

func change(n *int) {
	n = new(int) // возвращает новый указатель на int (o - zero value)
	// new(T) - Всегда возвращает указатель на новую переменную типа Т
	fmt.Println(n)
	*n = 20
}

func main() {
	x := 10

	change(&x)

	fmt.Println(x)
	stack := stack.New([]int{1, 2, 3})

	fmt.Println(*stack)
}
