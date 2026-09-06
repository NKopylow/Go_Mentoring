package main

import "fmt"

type MyType struct {
	a int
	b int
} // объявление типа через структуру + ключевое слово type

type Person struct {
	name string
}

// Задача на указатели
/*
1) Что выведет программа?
package main

import "fmt"

func change(n *int) {
	*n = 100
}

func main() {
	x := 10

	change(&x)

	fmt.Println(x)
}
*/

type User struct {
	Name string
	Age  int
}

func increaseAge(user *User) {
	user.Age++
}

/*
2) Напиши функцию: func increaseAge(user *User)
которая увеличивает возраст пользователя на 1.

Дан тип:
type User struct {
	Name string
	Age  int
}

И код:
func main() {
	user := User{
		Name: "Nikita",
		Age:  24,
	}

	increaseAge(&user)

	fmt.Println(user.Age)
}

Требование: после вызова должно быть: 25
*/

/*
Найти и исправить ошибку

package main

import "fmt"

func change(n *int) {
	n = new(int)
	*n = 20
}

func main() {
	x := 10

	change(&x)

	fmt.Println(x)
}
*/

// func main() {
// 	// var a bool = true
// 	// var b int = 10
// 	// var c float64 = 3.14
// 	// var d string = "Hello, World!"
// 	// var e [5]int = [5]int{1, 2, 3, 4, 5}
// 	// var f []int = []int{1, 2, 3, 4, 5}
// 	// var h *int = new(int)
// 	// var i func(int) int = func(int) int { return 1 }
// 	// // var j interface { ... } = interface { ... }
// 	// k := map[string]int{"one": 1, "two": 2, "three": 3}
// 	// // var l chan int = make(chan int)

// 	// for index, value := range k {
// 	// 	fmt.Println(index, value)
// 	// }

// 	// foo()

// 	user := User{
// 		Name: "Nikita",
// 		Age:  24,
// 	}

// 	increaseAge(&user)

// 	fmt.Println(user.Age)

// 	// fmt.Println(a, b, c, d, e, f, h, i, k, l)
// }

func foo() {
	fmt.Println("foo")
}

// Boolean — bool
// Numeric — целые, float, complex
// String — string
// Array — [5]int
// Slice — []int
// Struct — struct { ... }
// Pointer — *int
// Function — func(int) int
// Interface — interface { ... }
// Map — map[string]int
// Channel — chan int
