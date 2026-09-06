package main

// Напишите функцию GetOrCreate, которая создает новый элемент мапы,
// если его еще не было и возвращает его значение, или просто возвращает значение при наличии.
// Важно учесть, что код должен нормально работать в конкурентной среде.
import (
	"fmt"
	"sync"
)

type ConcurrentMap struct {
	rwMutex sync.RWMutex
	m       map[string]string
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{m: make(map[string]string)}
}

func (cm *ConcurrentMap) GetOrCreate(key string, value string) string {
	cm.rwMutex.Lock()
	defer cm.rwMutex.Unlock()
	_, ok := cm.m[key]

	if !ok {
		cm.m[key] = value
		return value
	} else {
		return value
	}
}

func main() {
	cm := NewConcurrentMap()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		val := cm.GetOrCreate("key1", "value1")
		fmt.Println("Goroutine 1 got:", val)
	}()

	go func() {
		defer wg.Done()
		val := cm.GetOrCreate("key1", "value2")
		fmt.Println("Goroutine 2 got:", val)
	}()
	wg.Wait()
}
