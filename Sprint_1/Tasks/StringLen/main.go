package main

import (
	"fmt"
	"runtime"
	"unicode/utf8"
	"unsafe"
)

func main() {
	var eng rune = 'G' // 4 байта
	var ru rune = 'П'  // 4 байта
	// ch := make(chan string)
	// go printSomeString(ch)
	// result := <-ch
	// fmt.Println(result)
	fmt.Println("--------------------------------")
	fmt.Println(runtime.NumGoroutine())
	fmt.Println("--------------------------------")
	fmt.Println("en: ", unsafe.Sizeof(eng))
	fmt.Println("ru: ", unsafe.Sizeof(ru))
	fmt.Println("--------------------------------")
	var str string = "ПOА"
	fmt.Println(string(str[0]))
	fmt.Println(rune(str[0]))
	fmt.Println("--------------------------------")
	for _, rune := range str { // range проходит по рунам
		fmt.Println(string(rune)) // string(rune) преобразует руну в строку, так выведутся конкретные символы
	}
	fmt.Println("--------------------------------")
	for i := 0; i < utf8.RuneCountInString(str); i++ {
		fmt.Println(str[i]) // str[i] проходит по байтам
	}
	for i := 0; i < len(str); i++ { // len возвращает количество байт в строке
		fmt.Println(str[i]) // str[i] проходит по байтам
	}

	fmt.Println("--------------------------------")
	a := []int{1, 2, 3}
	b := []int{23, 22, 25, 28}
	a = append(a, b...)
	fmt.Println("a: ", a)
	fmt.Println("--------------------------------")
	// fmt.Println(<-ch)
	struct1 := struct{}{}
	struct2 := struct{}{}
	fmt.Println(struct1 == struct2)
}

func printSomeString(ch chan<- string) {
	fmt.Println("Some string")
	ch <- "Some string"
}
