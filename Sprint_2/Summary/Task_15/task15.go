package main

import (
	"fmt"
	"strings"
)

// Тема: память и бенчмарки.
//
// Задача 15. Меньше аллокаций — быстрее код.
//
// Ниже FormatReport каждый раз создаёт много временных строк (Sprintf + +=).
// Это нагружает Garbage Collector.
//
// Шаги:
//   1. Добавь go.mod и task15_test.go с бенчмарком:
//
//        func BenchmarkFormatReport(b *testing.B) {
//            items := sampleItems()
//            b.ReportAllocs()
//            for b.Loop() {
//                _ = FormatReport(items)
//            }
//        }
//
//   2. Запусти: go test -bench=. -benchmem
//      Запиши в конец файла: ns/op, B/op, allocs/op.
//   3. Реализуй FormatReportBuilder — тот же текст, но через strings.Builder.
//   4. Добавь бенчмарк для Builder и сравни allocs/op.
//
// Подсказка для Builder:
//
//	var b strings.Builder
//	b.Grow(примерно нужный размер)  // необязательно, но полезно
//	fmt.Fprintf(&b, "#%d ...", ...)
//	return b.String()

type Item struct {
	ID    int
	Name  string
	Price int // рубли
	Qty   int
}

// FormatReport — наивная версия. Не меняй, это эталон для сравнения.
func FormatReport(items []Item) string {
	result := ""
	for _, it := range items {
		line := fmt.Sprintf("#%d %s x%d = %d", it.ID, it.Name, it.Qty, it.Price*it.Qty)
		result += line + "\n"
	}
	return result
}

// FormatReportBuilder должна вернуть ТОТ ЖЕ текст, что FormatReport,
// но собирать его через strings.Builder (без result += ...).
func FormatReportBuilder(items []Item) string {
	_ = strings.Builder{} // подсказка: используй этот тип
	panic("not implemented")
}

func sampleItems() []Item {
	return []Item{
		{ID: 1, Name: "Coffee", Price: 250, Qty: 2},
		{ID: 2, Name: "Bagel", Price: 120, Qty: 1},
		{ID: 3, Name: "Juice", Price: 180, Qty: 3},
	}
}

func main() {
	fmt.Print(FormatReport(sampleItems()))
	// После реализации раскомментируй:
	// fmt.Print(FormatReportBuilder(sampleItems()))
}

// Baseline FormatReport:     ns/op=???  B/op=???  allocs/op=???
// FormatReportBuilder:       ns/op=???  B/op=???  allocs/op=???
