package storage

import "coffee/internal/domain"

type MemoryStorage struct {
	stock   map[domain.Ingredient]int
	drinks  map[string]domain.Drink
	sales   []domain.Sale
	revenue int
}

func NewMemoryStorage(drinks []domain.Drink) *MemoryStorage {
	drinkMap := make(map[string]domain.Drink)

	for _, drink := range drinks {
		drinkMap[drink.Name] = drink
	}

	stock := make(map[domain.Ingredient]int)

	for _, ingredient := range domain.Ingredients {
		stock[ingredient] = 0
	}

	return &MemoryStorage{
		stock:  stock,
		drinks: drinkMap,
		sales:  []domain.Sale{},
	}
}

func (s *MemoryStorage) GetStock(ingredient domain.Ingredient) int {
	return s.stock[ingredient]
}

func (s *MemoryStorage) SetStock(
	ingredient domain.Ingredient,
	quantity int,
) {
	s.stock[ingredient] = quantity
}

func (s *MemoryStorage) AddStock(
	ingredient domain.Ingredient,
	quantity int,
) {
	s.stock[ingredient] += quantity
}

func (s *MemoryStorage) GetAllStock() map[domain.Ingredient]int {
	return s.stock
}

func (s *MemoryStorage) GetDrink(name string) (domain.Drink, bool) {
	drink, ok := s.drinks[name]
	return drink, ok
}

func (s *MemoryStorage) AddSale(sale domain.Sale) {
	s.sales = append(s.sales, sale)
	s.revenue += sale.Amount
}

func (s *MemoryStorage) GetStats() (int, int) {
	return len(s.sales), s.revenue
}
