package main

import "fmt"

// Что выведет программа и почему?

type Point struct {
	X, Y int
}

func move(p Point) { // значение-получатель
	p.X += 10
}

func main() {
	p1 := Point{1, 1}

	move(p1)

	fmt.Println(p1) // выведет {1,1} так как структура передается по значению
}
