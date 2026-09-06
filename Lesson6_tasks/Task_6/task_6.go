package main

import "fmt"

// Представьте, что у вас есть система, в которой нужно обрабатывать
// множество запросов, и для каждого из них требуется выполнить
// дорогостоящую операцию (например, запрос к базе данных
// или сложные вычисления).
//
// Чтобы не выполнять повторно одинаковую операцию для одного
// и того же ключа, вы хотите использовать кэш.
//
// Напишите функцию GetOrCompute, которая принимает:
//   - ключ;
//   - функцию для вычисления значения.
//
// Функция должна возвращать либо значение из кэша,
// либо результат вычисления.

func main() {
	cash := make(map[string]any)

	GetOrCompute("user:1", func() any {
		return "Nikita"
	}, cash)

	GetOrCompute("user:2", func() any {
		return 25
	}, cash)

	GetOrCompute("pi", func() any {
		return 3.14159
	}, cash)

	GetOrCompute("isAdmin", func() any {
		return true
	}, cash)

	GetOrCompute("expensive", func() any {
		return 42
	}, cash)

	GetOrCompute("expensive", func() any {
		return 100
	}, cash)
}

func GetOrCompute(key string, function func() any, cash map[string]any) any { // по условиям не указан тип возвращаемого значения функции в параметре

	if _, ok := cash[key]; ok {
		fmt.Printf("From cash ")
		return cash[key]
	} else {
		fmt.Printf("Computed value ")
		result := function()
		cash[key] = result
		return result
	}
}
