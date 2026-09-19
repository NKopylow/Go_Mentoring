package cli

import (
	"fmt"
	"strconv"

	"coffee/internal/domain"
	"coffee/internal/service"
)

const (
	exitOK       = 0
	exitError    = 1
	exitUsage    = 2
	exitInternal = 3
) // коды результата операции общения с кофеаппаратом

type CLI struct {
	service *service.CoffeeService
}

func New(service *service.CoffeeService) *CLI {
	return &CLI{
		service: service,
	}
}

func (c *CLI) Run(args []string) int {
	if len(args) == 0 {
		fmt.Println("usage")
		return exitUsage
	}

	switch args[0] { // тут некорректно считывается команда, надо смотреть всю строку, а не только первые 2 записи
	case "help":
		return c.help(args[1:])
	case "menu":
		return c.menu(args[1:])
	case "stock":
		return c.stock(args[1:])
	case "brew":
		return c.brew(args[1:])
	case "stats":
		return c.stats(args[1:])
	default:
		fmt.Println("usage")
		return exitUsage
	}
}

func (c *CLI) help(args []string) int {
	if len(args) != 1 {
		fmt.Println("usage")
		return exitUsage
	}

	fmt.Println("usage")

	return exitOK
}

func (c *CLI) menu(args []string) int {
	if len(args) != 1 {
		fmt.Println("usage")
		return exitUsage
	}

	for _, drink := range c.service.Menu() {
		fmt.Printf("%s %d\n", drink.Name, drink.Price) // %s — вывести строку (string)
		// %d — вывести целое число (int)
	}

	return exitOK
}

func (c *CLI) stock(args []string) int {
	if len(args) < 2 {
		fmt.Println("usage")
		return exitUsage
	}

	switch args[1] {
	case "get":
		return c.stockGet(args)

	case "add":
		return c.stockAdd(args)

	case "set":
		return c.stockSet(args)

	default:
		fmt.Println("usage")
		return exitUsage
	}
}

func (c *CLI) stockGet(args []string) int {
	if len(args) != 2 {
		fmt.Println("usage")
		return exitUsage
	}

	stock := c.service.GetStock()

	for _, ingredient := range domain.Ingredients {
		fmt.Printf(
			"%s=%d\n",
			ingredient,
			stock[ingredient],
		)
	}

	return exitOK
}

func (c *CLI) stockAdd(args []string) int {
	if len(args) != 4 {
		fmt.Println("usage")
		return exitUsage
	}

	quantity, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("некорректные параметры")
		return exitError
	}

	if err := c.service.AddStock(args[2], quantity); err != nil {
		fmt.Println(err)
		return exitError
	}

	fmt.Println("ok")

	return exitOK
}

func (c *CLI) stockSet(args []string) int {
	if len(args) != 4 {
		fmt.Println("usage")
		return exitUsage
	}

	quantity, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("некорректные параметры")
		return exitError
	}

	if err := c.service.SetStock(args[2], quantity); err != nil {
		fmt.Println(err)
		return exitError
	}

	fmt.Println("ok")

	return exitOK
}

func (c *CLI) brew(args []string) int {
	if len(args) != 4 || args[2] != "--pay" {
		fmt.Println("usage")
		return exitUsage
	}

	payment, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("некорректные параметры")
		return exitError
	}

	steps, err := c.service.Brew(args[1], payment)
	if err != nil {
		fmt.Println(err)
		return exitError
	}

	for _, step := range steps {
		fmt.Println(step)
	}

	return exitOK
}

func (c *CLI) stats(args []string) int {
	if len(args) != 1 {
		fmt.Println("usage")
		return exitUsage
	}

	orders, revenue := c.service.Stats()

	fmt.Printf(
		"orders=%d revenue=%d\n",
		orders,
		revenue,
	)

	return exitOK
}
