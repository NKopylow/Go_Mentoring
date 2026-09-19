package domain

type Drink struct {
	Name        string
	Price       int
	Ingredients map[Ingredient]int
	Steps       []string
}
