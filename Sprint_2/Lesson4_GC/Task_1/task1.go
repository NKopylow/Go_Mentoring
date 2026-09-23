package main

import "fmt"

func calculate() int {
	x := 42

	return x * 2
}

func main() {
	fmt.Println(calculate())
}
