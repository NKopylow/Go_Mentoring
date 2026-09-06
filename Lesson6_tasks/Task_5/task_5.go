package main

import (
	"fmt"
	"slices"
)

// У вас есть мапа, где ключом является строка, а значением — слайс строк.
// Каждый слайс содержит набор уникальных значений для этого ключа.
//
// Вам нужно написать функцию MergeToMap, которая принимает
// мапу и новый слайс для конкретного ключа, анализирует
// уже существующие значения в мапе и добавляет только
// те элементы из нового слайса, которых там ещё нет.
//
// Пример входных данных:
//
// Исходная мапа:
//
// m := map[string][]string{
//     "group1": {"apple", "banana"},
//     "group2": {"carrot"},
// }
//
// Новый слайс:
//
// newValues := []string{"banana", "cherry"}
//
// Ключ:
//
// key := "group1"
//
// Ожидаемый результат:
//
// m := map[string][]string{
//     "group1": {"apple", "banana", "cherry"},
//     "group2": {"carrot"},
// }

func MergeToMap(m map[string][]string, newValues []string, key string) map[string][]string {
	if newValues == nil || m == nil {
		return nil
	}

	currentMap := m[key] // слайс по ключу, который уже лежит в мапе

	for _, nValue := range newValues {
		if !slices.Contains(currentMap, nValue) {
			m[key] = append(m[key], nValue)
		}
	}
	return m
}

func main() {
	m := map[string][]string{
		"group1": {"apple", "banana"},
		"group2": {"carrot"},
	}

	newValues := []string{"banana", "cherry"}
	key := "group1"

	fmt.Println(MergeToMap(m, newValues, key))
}
