package main

import (
	"net/http"
	_ "net/http/pprof"
)

var storage [][]byte

func main() {
	go func() {
		for {
			storage = append(storage, make([]byte, 1024*1024))
		}
	}()

	http.ListenAndServe("localhost:6060", nil)
}
