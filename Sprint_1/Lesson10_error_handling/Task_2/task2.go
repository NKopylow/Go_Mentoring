package main

import "fmt"

func f() {
	defer fmt.Println("A")
	defer fmt.Println("B")
	fmt.Println("C")
}

func main() {
	f()
}

// реализуется логику FILO, поэтому будет С, В, А
