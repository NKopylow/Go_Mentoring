package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map
	var wg sync.WaitGroup

	const workers = 100

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()

			key := fmt.Sprintf("worker-%d", id)

			// Записываем
			m.Store(key, id)

			// Читаем
			value, ok := m.Load(key)

			if !ok {
				fmt.Println("value not found:", key)
				return
			}

			fmt.Println(key, "=", value)
		}(i)
	}

	wg.Wait()

	fmt.Println("All workers finished")
}
