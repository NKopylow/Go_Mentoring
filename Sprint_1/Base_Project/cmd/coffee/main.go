package main

import (
	"bufio"
	"os"
	"strings"

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

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		args := strings.Fields(scanner.Text())

		if len(args) == 0 {
			continue
		}

		if args[0] == "exit" {
			break
		}

		app.Run(args)
	}
}
