package main

import (
	"cart/internal/app"
	"log"
)

func main() {
	log.Println("Starting cart server")

	err := app.NewCartServiceApp()
	if err != nil {
		log.Printf("main.app.NewCartServiceApp: %+v", err.Error())
	}
}
