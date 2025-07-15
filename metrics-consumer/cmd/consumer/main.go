package main

import (
	"log"
	"metrics-consumer/internal"
)

func main() {
	err := internal.BootStrapMetricsService(".env")
	if err != nil {
		log.Printf("error building app bootstrapservice: %v\n", err.Error())
	}
}
