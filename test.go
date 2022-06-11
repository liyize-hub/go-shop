package main

import "go-shop/middleware"

func main() {
	rabbitmq := middleware.NewRabbitMQSimple("seckill")
	rabbitmq.ConsumeSimple()
}
