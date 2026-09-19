package main

import (
	"os"

	"coffee/internal/cli"
	"coffee/internal/domain"
	"coffee/internal/service"
	"coffee/internal/storage"
)

func main() {
	var drinksMemoryStorage []domain.Drink
	drinks := service.DefaultDrinks()

	for _, drink := range drinks {
		drinksMemoryStorage = append(drinksMemoryStorage, drink)
	}
	memoryStorage := storage.NewMemoryStorage(drinksMemoryStorage)

	coffeeService := service.NewCoffeeService(
		memoryStorage,
		drinksMemoryStorage,
	)

	app := cli.New(coffeeService)

	os.Exit(app.Run(os.Args[1:]))
}
