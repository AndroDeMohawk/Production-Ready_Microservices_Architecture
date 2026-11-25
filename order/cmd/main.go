package main

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/app"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load("./../order.env")
	app := app.New(context.Background())
	app.Run()
}
