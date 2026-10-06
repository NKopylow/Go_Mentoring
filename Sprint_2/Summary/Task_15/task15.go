package main

import (
	"fmt"
	"strings"
)

// Тема: управление памятью, GC, escape analysis, sync.Pool, бенчмарки.
//
// Задача 15. Снижаем аллокации в "горячем" пути.
//
// Ниже — наивная функция FormatReport, которая вызывается миллионы раз в секунду
// (например, при логировании). Она плохо работает с памятью.
//
// 1. Добавь go.mod и task15_test.go. Напиши бенчмарк BenchmarkFormatReport с b.ReportAllocs()
//    и зафиксируй baseline: ns/op, B/op, allocs/op (запиши комментарием в конце файла).
//
// 2. Выясни, ЧТО именно аллоцируется:
//    - go build -gcflags='-m -m' ./... — найди строки "escapes to heap" и "moved to heap";
//    - go test -bench . -memprofile mem.out && go tool pprof -alloc_space mem.out (top, list FormatReport).
//
// 3. Сделай оптимизированные версии, каждая — отдельная функция, чтобы сравнить в бенчмарке:
//    a) FormatReportBuilder   — strings.Builder с предварительным Grow;
//    b) FormatReportAppend    — работа с []byte через strconv.AppendInt/AppendQuote и
//                               возврат string(buf) (одна аллокация на результат);
//    c) FormatReportPool      — буфер берётся из sync.Pool (*bytes.Buffer или *[]byte) и
//                               возвращается обратно; подумай, почему в Pool лучше класть указатель,
//                               и почему слишком большой буфер стоит НЕ возвращать в Pool;
//    d) FormatReportTo(w io.Writer, ...) — вообще без возврата строки: пишем напрямую в writer.
//
// 4. Избавься от лишних escape'ов: передавай Item по указателю или по значению? Проверь оба
//    варианта gcflags'ом и бенчмарком — и объясни результат (размер структуры, инлайнинг).
//
// 5. Эксперимент с GC: напиши в main цикл, который вызывает наивную версию 5 млн раз, и запусти:
//      GODEBUG=gctrace=1 go run .      — посчитай, сколько циклов GC произошло;
//      GOGC=400 GODEBUG=gctrace=1 go run . — как изменилось количество циклов и RSS?
//    Затем то же для лучшей оптимизированной версии. Выводы — комментарием.
//
// 6. Добавь в бенчмарк проверку b.RunParallel для версии с Pool — именно под нагрузкой
//    из многих горутин Pool показывает преимущество. Объясни, почему.
//
// Цель: allocs/op у лучшей версии — 1 (или 0 для FormatReportTo).

type Item struct {
	ID    int64
	Name  string
	Price int64 // в копейках
	Qty   int32
	Tags  [4]string
}

// Наивная версия: много конкатенаций, fmt.Sprintf, промежуточные срезы.
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
//
// После оптимизаций:
//
// Выводы по GC (gctrace / GOGC):
//
