package main

import (
	"fmt"
	"slices"
)

// Вы разрабатываете сервис для подсчета уникальных слов,
// но количество слов может быть очень большим,
// // и вы решили ограничить максимальное число записе
// Если в мапу добавляется больше слов, чем указано лимитом,
// она должна автоматически удалять самые старые записи.

type WordCounter struct {
	counts map[string]int
	limit  int
	order  []string
}

func NewWordCounter(limit int) *WordCounter {
	return &WordCounter{
		counts: make(map[string]int),
		limit:  limit,
		order:  make([]string, 0),
	}
}

func (wc *WordCounter) CountWord(word string) {
	wc.counts[word]++
	if !slices.Contains(wc.order, word) {
		wc.order = append(wc.order, word)
	}

	if len(wc.counts) > wc.limit {
		// Логика удаления (здесь нужно реализовать)
		delete(wc.counts, wc.order[0])
		wc.order = wc.order[1:]
	}
}

func main() {
	wc := NewWordCounter(3)

	words := []string{
		"apple",  // должен быть удален
		"banana", // должен быть удален
		"apple",
		"orange",
		"grape",
		"banana",
		"kiwi",
	}

	for _, word := range words {
		wc.CountWord(word)
	}

	fmt.Println("Количество слов:", wc.counts)
	fmt.Println("Порядок слов:", wc.order)
}
