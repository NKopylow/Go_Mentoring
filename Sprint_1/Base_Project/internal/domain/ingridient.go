package domain

import "slices"

type Ingredient string

const (
	Beans Ingredient = "beans"
	Water Ingredient = "water"
	Milk  Ingredient = "milk"
	Sugar Ingredient = "sugar"
)

var Ingredients = []Ingredient{
	Beans,
	Water,
	Milk,
	Sugar,
}

func IsValidIngredient(ingredient Ingredient) bool {
	return slices.Contains(Ingredients, ingredient)
}
