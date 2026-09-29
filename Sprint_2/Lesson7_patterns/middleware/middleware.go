package main

import (
	"fmt"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("1. Logging")

		next.ServeHTTP(w, r)

		fmt.Println("5. Logging after")
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("2. Auth")

		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)

		fmt.Println("4. Auth after")
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("3. Handler")
	fmt.Fprintln(w, "Hello!")
}

func main() {
	var handler http.Handler = http.HandlerFunc(Handler)

	handler = Auth(handler)
	handler = Logging(handler)

	http.Handle("/", handler)

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
