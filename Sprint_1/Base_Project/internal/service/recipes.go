package service

import "coffee/internal/domain"

const (
	Espresso  string = "espresso"
	Americano string = "americano"
	Latte     string = "latte"
)

func DefaultDrinks() map[string]domain.Drink {
	return map[string]domain.Drink{
		Espresso: {
			Name:  "espresso",
			Price: 150,
			Ingredients: map[domain.Ingredient]int{
				domain.Beans: 8,
				domain.Water: 30,
			},
			Steps: []string{
				"grind",
				"tamp",
				"brew 25s",
			},
		},

		Americano: {
			Name:  "americano",
			Price: 180,
			Ingredients: map[domain.Ingredient]int{
				domain.Beans: 8,
				domain.Water: 120,
			},
			Steps: []string{
				"grind",
				"brew 25s",
				"add water",
			},
		},

		Latte: {
			Name:  "latte",
			Price: 220,
			Ingredients: map[domain.Ingredient]int{
				domain.Beans: 8,
				domain.Water: 30,
				domain.Milk:  150,
			},
			Steps: []string{
				"grind",
				"brew 25s",
				"steam milk",
				"mix",
			},
		},
	}
}
