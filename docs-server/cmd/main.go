package main

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/static"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := static.FS.ReadFile("swagger-ui.html")
		if err != nil {
			http.Error(w, "UI not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		data, err := static.FS.ReadFile("generated/swagger.swagger.json")
		if err != nil {
			http.Error(w, "Swagger JSON not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	port := ":8081" // Отдельный порт для документации
	slog.Info("Запуск сервера документации Swagger", "адрес", port)
	inventoryURL, _ := url.Parse("http://localhost:40041")
	inventoryProxy := httputil.NewSingleHostReverseProxy(inventoryURL)
	mux.Handle("/api/v1/inventory/", inventoryProxy)

	paymentURL, _ := url.Parse("http://localhost:40042")
	paymentProxy := httputil.NewSingleHostReverseProxy(paymentURL)
	mux.Handle("/api/v1/payment/", paymentProxy)
	if err := http.ListenAndServe(port, mux); err != nil {
		slog.Error("Ошибка запуска сервера документации", "error", err)
		os.Exit(1)
	}
	
}
