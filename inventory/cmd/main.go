package main

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("./../inventory.env")
	a := app.New(context.Background())
	a.Run()
}
