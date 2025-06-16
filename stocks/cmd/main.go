package main

import (
	"log"
	"stocks/internal/app"
)

func main() {
	log.Println("Starting stocks server")

	err := app.NewStockServiceApp()
	if err != nil {
		log.Printf("main.app.NewStockServiceApp: %+v", err.Error())
	}
}
