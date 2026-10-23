package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	dsn := "host=localhost port=5432 user=postgres password=12345 dbname=join_practice sslmode=disable"

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("не удалось подключиться:", err)
	}
	defer conn.Close(ctx)

	if err := conn.Ping(ctx); err != nil {
		log.Fatal("ping не удался:", err)
	}

	const query = `
		SELECT users.name, COUNT(orders.order_number)
		FROM orders
		RIGHT JOIN users ON users.id = orders.user_id
		GROUP BY users.name
	`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var count int64
		if err := rows.Scan(&name, &count); err != nil {
			log.Fatal(err)
		}
		fmt.Println(name, count)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
