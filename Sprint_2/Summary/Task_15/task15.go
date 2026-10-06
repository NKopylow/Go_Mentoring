package main

import (
	"fmt"
	"strings"
)

// Тема: память, GC, escape analysis, sync.Pool, бенчмарки.
//
// Задача 15. Снижаем аллокации.
//
// Ниже — наивная FormatReport. Сделай:
//   1. Бенчмарк с b.ReportAllocs, зафиксируй baseline (ns/op, B/op, allocs/op) комментарием.
//   2. go build -gcflags='-m' — посмотри, что escapes to heap.
//   3. Две оптимизированные версии (см. ниже). Цель: ≤1 allocs/op у лучшей.
//   4. GODEBUG=gctrace=1 go run . для наивной и лучшей — сколько циклов GC? Краткий вывод.

type Item struct {
	ID    int64
	Name  string
	Price int64 // копейки
	Qty   int32
	Tags  [4]string
}

// FormatReport — наивная версия (много Sprintf и конкатенаций). Не меняй, это baseline.
func FormatReport(items []Item) string {
	result := ""
	for _, it := range items {
		line := fmt.Sprintf("#%d %q x%d = %d.%02d", it.ID, it.Name, it.Qty,
			it.Price*int64(it.Qty)/100, it.Price*int64(it.Qty)%100)
		tags := []string{}
		for _, t := range it.Tags {
			if t != "" {
				tags = append(tags, t)
			}
		}
		if len(tags) > 0 {
			line += " [" + strings.Join(tags, ",") + "]"
		}
		result += line + "\n"
	}
	return result
}

// FormatReportBuilder собирает тот же текст через strings.Builder (с Grow заранее).
// Результат должен совпадать с FormatReport.
func FormatReportBuilder(items []Item) string {
	panic("not implemented")
}

// FormatReportPool берёт буфер из sync.Pool, пишет туда, возвращает string, буфер — обратно в Pool.
// Результат должен совпадать с FormatReport.
func FormatReportPool(items []Item) string {
	panic("not implemented")
}

func sampleItems() []Item {
	return []Item{
		{ID: 1, Name: "Coffee", Price: 25000, Qty: 2, Tags: [4]string{"hot", "drink"}},
		{ID: 2, Name: "Bagel", Price: 12050, Qty: 1},
		{ID: 3, Name: "Juice", Price: 18000, Qty: 3, Tags: [4]string{"cold", "drink", "fresh"}},
	}
}

func main() {
	fmt.Print(FormatReport(sampleItems()))
}

// Baseline:
// После оптимизаций:
// Выводы по GC:
