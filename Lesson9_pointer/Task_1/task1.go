package main

import "fmt"

// TODO: допиши функцию inc//
// должна увеличить значение по адресу на 1

func inc(x *int) {
	*x = *x + 1
}

func main() {
	x := 10
	inc(&x)
	fmt.Println(x)
} // должно вывести 11}
