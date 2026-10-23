-- Выполни INNER JOIN двух таблиц users и orders, 
-- выбрав имя пользователя и номер заказа.

-- SELECT name, order_number from users
-- inner join orders on users.id = orders.user_id

-- Напиши запрос с LEFT JOIN, выводящий всех пользователей 
-- и, если есть, их заказы.

-- SELECT name, order_number FROM users
-- LEFT JOIN orders on users.id = user_id

-- Используй GROUP BY и COUNT() для подсчета количества заказов 
-- по каждому пользователю.

-- SELECT name, COUNT(order_number)
-- FROM orders
-- RIGHT JOIN users on users.id = user_id 
-- GROUP BY users.name


-- Напиши подзапрос в WHERE, 
-- чтобы найти пользователей, у которых есть заказы.

-- SELECT name 
-- FROM users
-- WHERE id in (SELECT users.id FROM users RIGHT JOIN orders on user_id = users.id)

-- Выполни CROSS JOIN таблиц colors и sizes — создай декартово 
-- произведение комбинаций.

SELECT *
FROM orders
CROSS JOIN users




