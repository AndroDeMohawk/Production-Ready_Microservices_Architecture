package main

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("./../payment.env")
	a := app.New(context.Background())
	a.Run()

}
