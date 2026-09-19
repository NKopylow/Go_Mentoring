package main

import (
	"fmt"
)

// Напиши программу, которая поменяет местами значения двух переменных.

// TODO: func swap(a, b *int)

func swap(a, b *int) {
	*a, *b = *b, *a
}

func main() {
	a, b := 3, 5

	swap(&a, &b)

	fmt.Println("a=", a, "b=", b)
}
