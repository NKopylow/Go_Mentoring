package main

import (
	"fmt"
	"sync"
)

// Покажи, как с помощью sync.Once инициализировать структуру только один раз.

type Person struct {
	id  int
	age int
}

func main() {
	var once sync.Once
	var wg sync.WaitGroup
	p := Person{}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(func() {
				p = Person{
					id:  i,
					age: i * 10,
				}
				fmt.Println("p: ", p)
			})
		}()
	}
	wg.Wait()
	fmt.Println(p)
}
