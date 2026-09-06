package task3

import (
	"fmt"
	"time"
)

// Что выведет программа и почему?
// Также проанализируй код и исправь все потенциальные проблемы, если таковые есть.

func UpdateProductStock() <-chan map[string]int {
	stockUpdates := make(chan map[string]int) // новый канал
	go func() {
		currentStock := map[string]int{
			"Apples":  50,
			"Bananes": 30,
			"Oranges": 20,
			"Grapes":  15,
		}
		for i := 0; i < 5; i++ {
			for product, quantity := range currentStock {
				println("quantity: ", quantity, "int(float64(quantity) * 0.95): ", int(float64(quantity)*0.95))
				currentStock[product] = int(float64(quantity) * 0.95) // добавляет скидку 5%
			}
			stockUpdates <- currentStock
			time.Sleep(150 * time.Millisecond) // добавляет задержку 150 миллисекунд
		}
	}()
	return stockUpdates
}

func main() {
	stockStream := UpdateProductStock()
	var stockHistory []map[string]int
	for i := 0; i < 5; i++ {
		stock := <-stockStream                     // постепенно получает новые стоимости продуктов
		stockHistory = append(stockHistory, stock) // добавляет новые стоимости продуктов в историю
	}
	for i, stock := range stockHistory {
		fmt.Printf("Iteration %d: %v\n", i+1, stock)
	}
} // 5 итераций
// 1 2 3 4 5

// stock - каждый раз стоимость продуктов со скидкой 5% от предыдущей стоимости
