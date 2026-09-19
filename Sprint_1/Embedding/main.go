package main

import "fmt"

type Person struct {
	name string
	age  int8
}

func (p *Person) ChangeName(name string) {
	fmt.Println("Person.ChangeName")
	p.name = name
}

func (p *Person) GetName() string {
	return p.name
}

type Employee struct {
	Person
	salary int32
}

func (e *Employee) ChangeName(name string) {
	fmt.Println("Employee.ChangeName")
	e.Person.name = name
}

func main() {
	e := Employee{name: "Anna", age: 20, salary: 1000}

	e.Person.ChangeName("John")

	fmt.Println(e.GetName())
}

type Flyer interface {
	Fly()
}

type Swimmer interface {
	Swim()
}

type Duck interface {
	Flyer
	Swimmer
}
