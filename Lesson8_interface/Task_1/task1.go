package main

import "fmt"

// Что сделать: запусти программу, убедись, что она выводит два разных сообщения, хотя функция получает интерфейс.

type Greeter interface {
	Greet()
}

type EN struct{}

func (EN) Greet() {
	fmt.Println("Hello!")
}

type RU struct{}

func (RU) Greet() {
	fmt.Println("Привет!")
}

func Say(g Greeter) {
	g.Greet()
}

func main() {
	Say(EN{})
	Say(RU{})
}
