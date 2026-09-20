package main

import (
	"fmt"
)

// Напиши пример с буферизированным каналом,
// где отправка не блокируется до тех пор, пока буфер не заполнится.

// full - сигнал заполненности канала, пустая структура, потому что при отправке в такой канал сообщение не будет ничего весить
func writer(channel chan int) {
	channel <- 0
	channel <- 1

	fmt.Println("buffer is full")

	channel <- 100

	fmt.Println("100 sent")
}

func main() {
	channel := make(chan int, 2)

	go writer(channel)

	fmt.Println(<-channel)
	fmt.Println(<-channel)
	fmt.Println(<-channel)
}
