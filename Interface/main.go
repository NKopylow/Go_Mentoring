package main

import "fmt"

type Duck interface {
	Quack()
	Fly()
	Swim()
	Walk()
}

// class SomeDuck implements Duck {} - в go такого нет

type EstairesDuck struct{}

type Mallard struct{}

func (m *Mallard) Fly() {
	fmt.Println("Mallard is flying")
}
func (m *Mallard) Quack() {}
func (m *Mallard) Swim()  {}
func (m *Mallard) Walk()  {}

func processDuck(d Duck) {
	d.Fly()
}

func (e *EstairesDuck) Quack() {}
func (e *EstairesDuck) Fly() {
	fmt.Println("EstairesDuck is flying")
}
func (e *EstairesDuck) Swim() {}
func (e *EstairesDuck) Walk() {}

func main() {
	e := EstairesDuck{}
	processDuck(&e)

	m := Mallard{}
	processDuck(&m)
}
